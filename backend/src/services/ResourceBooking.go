package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

type ResourceBookingService struct {
	db     *gorm.DB
	repo   *repositories.ResourceBookingRepository
	audit  *repositories.AuditLogRepository
	locker *utils.TurnaroundLocker
}

func NewResourceBookingService(db *gorm.DB, repo *repositories.ResourceBookingRepository, audit *repositories.AuditLogRepository, locker *utils.TurnaroundLocker) *ResourceBookingService {
	return &ResourceBookingService{db: db, repo: repo, audit: audit, locker: locker}
}

func (s *ResourceBookingService) List() ([]types.ResourceBookingDTO, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	return constructors.NewResourceBookingListResponse(rows), nil
}

// ResolveConflict 解决预约冲突：CONFLICT -> ACTIVE。ACTIVE 预约只在放行事务里一次性释放。
// 先按航班加锁再开事务，避免与放行事务互锁。
func (s *ResourceBookingService) ResolveConflict(id int64, actor string) (*types.ResourceBookingDTO, error) {
	var preload models.ResourceBooking
	if err := s.db.Preload("Resource").First(&preload, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.NewAppError(constants.BookingNotFound, fmt.Sprintf(constants.BookingNotFoundMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	if preload.BookingStatus != constants.BookingConflict {
		// 非冲突预约无需解决，幂等返回当前态。
		dto := constructors.NewResourceBookingResponse(preload)
		return &dto, nil
	}

	unlock := s.locker.Lock(preload.TurnaroundID)
	defer unlock()

	var booking models.ResourceBooking
	err := s.db.Transaction(func(tx *gorm.DB) error {
		b, err := s.repo.GetByID(tx, id)
		if err != nil {
			return err
		}
		booking = *b
		affected, err := s.repo.ResolveConflict(tx, id)
		if err != nil {
			return err
		}
		if affected == 0 {
			return gorm.ErrRecordNotFound
		}

		detail := fmt.Sprintf(constants.LogTemplates["ResourceBooking"][2], id, constants.BookingConflict, constants.BookingActive)
		return s.audit.Append(tx, models.AuditLog{
			Actor: actor, Action: "BOOKING_RESOLVE_CONFLICT", TargetType: "ResourceBooking",
			TargetID: fmt.Sprintf("%d", id), Detail: detail, CreatedAt: time.Now(),
		})
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.NewAppError(constants.BookingNotFound, fmt.Sprintf(constants.BookingNotFoundMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	booking.BookingStatus = constants.BookingActive
	booking.ConflictReason = ""
	dto := constructors.NewResourceBookingResponse(booking)
	return &dto, nil
}

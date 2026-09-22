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

type DelayEventService struct {
	db     *gorm.DB
	repo   *repositories.DelayEventRepository
	audit  *repositories.AuditLogRepository
	locker *utils.TurnaroundLocker
}

func NewDelayEventService(db *gorm.DB, repo *repositories.DelayEventRepository, audit *repositories.AuditLogRepository, locker *utils.TurnaroundLocker) *DelayEventService {
	return &DelayEventService{db: db, repo: repo, audit: audit, locker: locker}
}

func (s *DelayEventService) List() ([]types.DelayEventDTO, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	return constructors.NewDelayEventListResponse(rows), nil
}

// Resolve 关闭延误：未关闭延误计数归零后，航班才有资格放行。
// 先按航班加锁再开事务，避免与放行事务互锁。
func (s *DelayEventService) Resolve(id int64, actor string) (*types.DelayEventDTO, error) {
	var preload models.DelayEvent
	if err := s.db.First(&preload, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.NewAppError(constants.DelayNotFound, fmt.Sprintf(constants.DelayNotFoundMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	if preload.ResolvedAt != nil {
		// 已关闭视为幂等成功，直接返回当前态。
		dto := constructors.NewDelayEventResponse(preload)
		return &dto, nil
	}

	unlock := s.locker.Lock(preload.TurnaroundID)
	defer unlock()

	var delay models.DelayEvent
	err := s.db.Transaction(func(tx *gorm.DB) error {
		d, err := s.repo.GetByID(tx, id)
		if err != nil {
			return err
		}
		delay = *d
		now := time.Now()
		affected, err := s.repo.Resolve(tx, id, now)
		if err != nil {
			return err
		}
		if affected == 0 {
			return gorm.ErrRecordNotFound
		}

		detail := fmt.Sprintf(constants.LogTemplates["DelayEvent"][3], id, delay.TurnaroundID)
		return s.audit.Append(tx, models.AuditLog{
			Actor: actor, Action: "DELAY_RESOLVE", TargetType: "DelayEvent",
			TargetID: fmt.Sprintf("%d", id), Detail: detail, CreatedAt: now,
		})
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.NewAppError(constants.DelayNotFound, fmt.Sprintf(constants.DelayNotFoundMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	now := time.Now()
	delay.ResolvedAt = &now
	dto := constructors.NewDelayEventResponse(delay)
	return &dto, nil
}

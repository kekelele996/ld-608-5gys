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

// FlightTurnaroundService 承载过站放行联动事务。
type FlightTurnaroundService struct {
	db       *gorm.DB
	repo     *repositories.FlightTurnaroundRepository
	taskRepo *repositories.GroundTaskRepository
	bookRepo *repositories.ResourceBookingRepository
	audit    *repositories.AuditLogRepository
	locker   *utils.TurnaroundLocker
}

func NewFlightTurnaroundService(
	db *gorm.DB,
	repo *repositories.FlightTurnaroundRepository,
	taskRepo *repositories.GroundTaskRepository,
	bookRepo *repositories.ResourceBookingRepository,
	audit *repositories.AuditLogRepository,
	locker *utils.TurnaroundLocker,
) *FlightTurnaroundService {
	return &FlightTurnaroundService{db: db, repo: repo, taskRepo: taskRepo, bookRepo: bookRepo, audit: audit, locker: locker}
}

func (s *FlightTurnaroundService) List() ([]types.FlightTurnaroundDTO, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	out := make([]types.FlightTurnaroundDTO, 0, len(rows))
	for i := range rows {
		row := &rows[i]
		enriched, err := s.repo.GetWithAssociations(row.ID)
		if err != nil {
			continue
		}
		out = append(out, constructors.NewFlightTurnaroundResponse(*enriched))
	}
	return out, nil
}

func (s *FlightTurnaroundService) Get(id int64) (*types.FlightTurnaroundDTO, error) {
	row, err := s.repo.GetWithAssociations(id)
	if err != nil {
		return nil, types.NewAppError(constants.TurnaroundNotFound,
			fmt.Sprintf(constants.TurnaroundNotFoundMessage, id))
	}
	dto := constructors.NewFlightTurnaroundResponse(*row)
	return &dto, nil
}

// Gate 返回某航班的放行门禁快照（任务/延误/预约阻塞清单）。
func (s *FlightTurnaroundService) Gate(id int64) (*types.ReleaseGateDTO, error) {
	row, err := s.repo.GetWithAssociations(id)
	if err != nil {
		return nil, types.NewAppError(constants.TurnaroundNotFound,
			fmt.Sprintf(constants.TurnaroundNotFoundMessage, id))
	}
	gate := constructors.BuildGate(*row, row.Tasks, row.Delays, row.Bookings)
	return &gate, nil
}

// Release 过站放行：缺任一项就整次拒绝；通过后单事务进入 READY、释放 ACTIVE 预约、刷新任务与时间线。
func (s *FlightTurnaroundService) Release(id int64, actor string) (*types.ReleaseResultDTO, error) {
	unlock := s.locker.Lock(id)
	defer unlock()

	// 1) 聚合检查（事务内，避免脏读）。
	var gate types.ReleaseGateDTO
	var current *models.FlightTurnaround
	err := s.db.Transaction(func(tx *gorm.DB) error {
		row, err := s.repo.GetByID(tx, id)
		if err != nil {
			return err
		}
		current = row

		tasks, err := s.taskRepo.ListByTurnaround(tx, id)
		if err != nil {
			return err
		}
		var delays []models.DelayEvent
		if err := tx.Where("turnaround_id = ?", id).Order("id asc").Find(&delays).Error; err != nil {
			return err
		}
		var bookings []models.ResourceBooking
		if err := tx.Preload("Resource").Where("turnaround_id = ?", id).Order("id asc").Find(&bookings).Error; err != nil {
			return err
		}
		current.Tasks = tasks
		current.Delays = delays
		current.Bookings = bookings
		gate = constructors.BuildGate(*current, tasks, delays, bookings)
		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.NewAppError(constants.TurnaroundNotFound,
				fmt.Sprintf(constants.TurnaroundNotFoundMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}

	// 2) 重复放行：已 READY/DEPARTED 只提示，不再写任何数据。
	if constants.ReadyTerminalStatuses[current.TurnaroundStatus] {
		return nil, types.NewAppError(constants.ReleaseAlreadyDone,
			fmt.Sprintf(constants.ReleaseAlreadyDoneMessage, current.FlightNo))
	}

	// 3) 门禁不通过：整次拒绝，航班、任务、预约保持原样，分别返回任务、延误、预约阻塞清单。
	if !gate.Releasable {
		return nil, types.NewAppError(constants.ReleaseBlocked,
			fmt.Sprintf(constants.ReleaseBlockedMessage, current.FlightNo)).WithDetails(map[string]any{
			"turnaround_id":     id,
			"unsigned_tasks":    gate.UnsignedTasks,
			"open_delays":       gate.OpenDelays,
			"blocking_bookings": gate.BlockingBookings,
		})
	}

	// 4) 门禁通过：单事务完成状态推进 + ACTIVE 预约一次性释放 + 时间线同步。
	now := time.Now()
	var releasedCount int64
	err = s.db.Transaction(func(tx *gorm.DB) error {
		affected, err := s.repo.MarkReady(tx, id, now)
		if err != nil {
			return err
		}
		// 条件更新 0 行：并发提交中已有一次放行生效。
		if affected == 0 {
			return errReleaseRace
		}

		bookingIDs, err := s.bookRepo.ReleaseActiveForTurnaround(tx, id, now)
		if err != nil {
			return err
		}
		releasedCount = int64(len(bookingIDs))
		if releasedCount > 0 {
			resourceIDs, err := s.bookRepo.ListResourcesByBookingIDs(tx, bookingIDs)
			if err != nil {
				return err
			}
			if err := s.bookRepo.MarkResourcesAvailable(tx, resourceIDs); err != nil {
				return err
			}
		}

		if err := s.taskRepo.TouchOnRelease(tx, id, now); err != nil {
			return err
		}

		logText := fmt.Sprintf(constants.LogTemplates["FlightTurnaround"][3],
			id, current.FlightNo, releasedCount, len(current.Tasks))
		return s.audit.Append(tx, models.AuditLog{
			Actor:      actor,
			Action:     "TURNAROUND_RELEASE",
			TargetType: "FlightTurnaround",
			TargetID:   fmt.Sprintf("%d", id),
			Detail:     logText,
			CreatedAt:  now,
		})
	})
	if err != nil {
		if err == errReleaseRace {
			return nil, types.NewAppError(constants.ReleaseRaceLost,
				fmt.Sprintf(constants.ReleaseRaceLostMessage, current.FlightNo))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}

	return &types.ReleaseResultDTO{
		TurnaroundID:     id,
		FlightNo:         current.FlightNo,
		TurnaroundStatus: constants.TurnaroundReady,
		ReadyAt:          now.Format(time.RFC3339),
		ReleasedBookings: releasedCount,
		Gate:             gate,
	}, nil
}

var errReleaseRace = fmt.Errorf("release race lost")

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

type GroundTaskService struct {
	db       *gorm.DB
	repo     *repositories.GroundTaskRepository
	turnRepo *repositories.FlightTurnaroundRepository
	audit    *repositories.AuditLogRepository
	locker   *utils.TurnaroundLocker
}

func NewGroundTaskService(db *gorm.DB, repo *repositories.GroundTaskRepository, turnRepo *repositories.FlightTurnaroundRepository, audit *repositories.AuditLogRepository, locker *utils.TurnaroundLocker) *GroundTaskService {
	return &GroundTaskService{db: db, repo: repo, turnRepo: turnRepo, audit: audit, locker: locker}
}

func (s *GroundTaskService) List() ([]types.GroundTaskDTO, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	return constructors.NewGroundTaskListResponse(rows), nil
}

// Sign 签收任务：SIGNED 后与过站时间线同步刷新（actual_finish 落值，审计日志记录）。
// 先按航班加锁再开事务，避免与放行事务在单连接数据库上互锁。
func (s *GroundTaskService) Sign(id int64, actor string) (*types.GroundTaskDTO, error) {
	var preload models.GroundTask
	if err := s.db.First(&preload, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.NewAppError(constants.TaskNotFound, fmt.Sprintf(constants.TaskNotFoundMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	if preload.Status == constants.TaskSigned {
		return nil, types.NewAppError(constants.TaskAlreadySigned, fmt.Sprintf(constants.TaskAlreadySignedMessage, id))
	}

	unlock := s.locker.Lock(preload.TurnaroundID)
	defer unlock()

	var task models.GroundTask
	err := s.db.Transaction(func(tx *gorm.DB) error {
		t, err := s.repo.GetByID(tx, id)
		if err != nil {
			return err
		}
		task = *t
		now := time.Now()
		affected, err := s.repo.Sign(tx, id, now)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errAlreadySigned
		}

		detail := fmt.Sprintf(constants.LogTemplates["GroundTask"][3], id, task.TaskType, task.TurnaroundID)
		return s.audit.Append(tx, models.AuditLog{
			Actor: actor, Action: "TASK_SIGN", TargetType: "GroundTask",
			TargetID: fmt.Sprintf("%d", id), Detail: detail, CreatedAt: now,
		})
	})
	if err != nil {
		if err == errAlreadySigned {
			return nil, types.NewAppError(constants.TaskAlreadySigned, fmt.Sprintf(constants.TaskAlreadySignedMessage, id))
		}
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}

	now := time.Now()
	task.Status = constants.TaskSigned
	task.ActualFinish = &now
	dto := constructors.NewGroundTaskResponse(task)
	return &dto, nil
}

var errAlreadySigned = fmt.Errorf("task already signed")

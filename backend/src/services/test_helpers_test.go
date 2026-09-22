package services

import (
	"gorm.io/gorm"

	"groundTurn/src/repositories"
	"groundTurn/src/utils"
)

func newTaskService(db *gorm.DB) *GroundTaskService {
	locker := utils.NewTurnaroundLocker()
	return NewGroundTaskService(
		db,
		repositories.NewGroundTaskRepository(db),
		repositories.NewFlightTurnaroundRepository(db),
		repositories.NewAuditLogRepository(db),
		locker,
	)
}

func newDelayService(db *gorm.DB) *DelayEventService {
	return NewDelayEventService(db, repositories.NewDelayEventRepository(db), repositories.NewAuditLogRepository(db), utils.NewTurnaroundLocker())
}

func newBookingService(db *gorm.DB) *ResourceBookingService {
	return NewResourceBookingService(db, repositories.NewResourceBookingRepository(db), repositories.NewAuditLogRepository(db), utils.NewTurnaroundLocker())
}

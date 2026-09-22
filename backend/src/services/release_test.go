package services

import (
	"sync"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/utils"
)

func newTestService(t *testing.T) (*FlightTurnaroundService, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(
		&models.FlightTurnaround{},
		&models.GroundTask{},
		&models.GroundResource{},
		&models.ResourceBooking{},
		&models.DelayEvent{},
		&models.AuditLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	config.SeedIfEmpty(db)

	locker := utils.NewTurnaroundLocker()
	turnRepo := repositories.NewFlightTurnaroundRepository(db)
	taskRepo := repositories.NewGroundTaskRepository(db)
	bookRepo := repositories.NewResourceBookingRepository(db)
	auditRepo := repositories.NewAuditLogRepository(db)
	svc := NewFlightTurnaroundService(db, turnRepo, taskRepo, bookRepo, auditRepo, locker)
	return svc, db
}

func TestRelease_BlockedKeepsEverythingIntact(t *testing.T) {
	svc, db := newTestService(t)

	// TA2（id=2）三类阻塞齐全：放行必须失败。
	_, err := svc.Release(2, "tester")
	if err == nil {
		t.Fatal("expected release to be blocked, got nil error")
	}
	appErr := err.Error()
	if appErr == "" {
		t.Fatal("expected non-empty error")
	}

	// 航班仍 IN_SERVICE。
	var ta models.FlightTurnaround
	if err := db.First(&ta, 2).Error; err != nil {
		t.Fatal(err)
	}
	if ta.TurnaroundStatus != constants.TurnaroundInService {
		t.Fatalf("turnaround status changed on failed release: %s", ta.TurnaroundStatus)
	}
	if ta.ReadyAt != nil {
		t.Fatal("ready_at should stay nil on failed release")
	}

	// 任务、预约均保持原样。
	var statusCount int64
	db.Model(&models.GroundTask{}).Where("turnaround_id = ? AND status <> ?", 2, constants.TaskSigned).Count(&statusCount)
	if statusCount != 3 {
		t.Fatalf("expected 3 unsigned tasks unchanged, got %d", statusCount)
	}
	var conflictCount int64
	db.Model(&models.ResourceBooking{}).Where("turnaround_id = ? AND booking_status = ?", 2, constants.BookingConflict).Count(&conflictCount)
	if conflictCount != 1 {
		t.Fatalf("conflict booking must remain CONFLICT, got %d", conflictCount)
	}
}

func TestRelease_GreenFlowReleasesActiveBookings(t *testing.T) {
	svc, db := newTestService(t)

	// TA1（id=1）门禁全绿：放行成功并一次性释放 3 个 ACTIVE 预约。
	result, err := svc.Release(1, "tester")
	if err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if result.TurnaroundStatus != constants.TurnaroundReady {
		t.Fatalf("status = %s, want READY", result.TurnaroundStatus)
	}
	if result.ReleasedBookings != 3 {
		t.Fatalf("released bookings = %d, want 3", result.ReleasedBookings)
	}

	var ta models.FlightTurnaround
	if err := db.First(&ta, 1).Error; err != nil {
		t.Fatal(err)
	}
	if ta.ReadyAt == nil {
		t.Fatal("ready_at should be set")
	}

	var remaining int64
	db.Model(&models.ResourceBooking{}).
		Where("turnaround_id = ? AND booking_status = ?", 1, constants.BookingActive).Count(&remaining)
	if remaining != 0 {
		t.Fatalf("active bookings remain: %d", remaining)
	}

	var released int64
	db.Model(&models.ResourceBooking{}).
		Where("turnaround_id = ? AND booking_status = ?", 1, constants.BookingReleased).Count(&released)
	if released != 3 {
		t.Fatalf("released = %d, want 3", released)
	}

	// 资源回到 AVAILABLE。
	var available int64
	db.Model(&models.GroundResource{}).Where("availability_status = ?", constants.ResourceAvailable).Count(&available)
	if available != 3 {
		t.Fatalf("available resources = %d, want 3", available)
	}
}

func TestRelease_DuplicateAndConcurrentOnlyOnce(t *testing.T) {
	svc, db := newTestService(t)

	// 顺序重复放行：第二次必须失败（RELEASE_ALREADY_DONE）。
	if _, err := svc.Release(1, "tester"); err != nil {
		t.Fatalf("first release failed: %v", err)
	}
	if _, err := svc.Release(1, "tester"); err == nil {
		t.Fatal("duplicate release should fail")
	}

	// 重置数据库场景，验证并发提交精确一次。
	svc2, db2 := newTestService(t)
	_ = db
	var wg sync.WaitGroup
	const n = 20
	success := make(chan struct{}, n)
	failed := make(chan struct{}, n)
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if _, err := svc2.Release(1, "tester"); err == nil {
				success <- struct{}{}
			} else {
				failed <- struct{}{}
			}
		}()
	}
	close(start)
	wg.Wait()
	if len(success) != 1 {
		t.Fatalf("concurrent releases: success = %d, want exactly 1", len(success))
	}
	if len(failed) != n-1 {
		t.Fatalf("concurrent releases: failed = %d, want %d", len(failed), n-1)
	}

	var readyCount int64
	db2.Model(&models.FlightTurnaround{}).Where("id = 1 AND turnaround_status = ?", constants.TurnaroundReady).Count(&readyCount)
	if readyCount != 1 {
		t.Fatalf("ready rows = %d, want 1", readyCount)
	}
}

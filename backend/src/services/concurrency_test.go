package services

import (
	"sync"
	"testing"

	"groundTurn/src/constants"
	"groundTurn/src/models"
)

// TestConcurrentSignThenRelease 同一航班上并发签收任务与放行，最终状态必须一致：
// 要么放行先发生（被阻塞拒绝，签收随后成功），要么签收全部完成后放行成功 —— 不能死锁、不能半成品。
func TestConcurrentSignThenRelease(t *testing.T) {
	svc, db := newTestService(t)

	// TA2（id=2）有 3 个未签收任务、1 个未关闭延误、1 个 CONFLICT 预约。
	taskSvc := newTaskService(db)
	delaySvc := newDelayService(db)
	bookingSvc := newBookingService(db)

	var wg sync.WaitGroup
	actions := []func(){
		func() { _, _ = taskSvc.Sign(4, "t") },
		func() { _, _ = taskSvc.Sign(5, "t") },
		func() { _, _ = taskSvc.Sign(6, "t") },
		func() { _, _ = delaySvc.Resolve(1, "t") },
		func() { _, _ = bookingSvc.ResolveConflict(5, "t") },
		func() { _, _ = svc.Release(2, "t") }, // 极可能先于解阻完成而被阻塞
	}
	start := make(chan struct{})
	for _, action := range actions {
		wg.Add(1)
		go func(fn func()) {
			defer wg.Done()
			<-start
			fn()
		}(action)
	}
	close(start)
	wg.Wait()

	// 再放一次行（此时所有阻塞必然已清空），必须成功。
	result, err := svc.Release(2, "t")
	if err != nil {
		t.Fatalf("release after all resolved failed: %v", err)
	}
	if result.TurnaroundStatus != constants.TurnaroundReady {
		t.Fatalf("status = %s, want READY", result.TurnaroundStatus)
	}
	if result.ReleasedBookings != 2 {
		t.Fatalf("released = %d, want 2 (booking4 + resolved booking5)", result.ReleasedBookings)
	}

	// 校验全部任务 SIGNED、延误已关闭、预约 RELEASED。
	var unsigned int64
	db.Model(&models.GroundTask{}).Where("turnaround_id = ? AND status <> ?", 2, constants.TaskSigned).Count(&unsigned)
	if unsigned != 0 {
		t.Fatalf("unsigned tasks = %d, want 0", unsigned)
	}
	var openDelays int64
	db.Model(&models.DelayEvent{}).Where("turnaround_id = ? AND resolved_at IS NULL", 2).Count(&openDelays)
	if openDelays != 0 {
		t.Fatalf("open delays = %d, want 0", openDelays)
	}
	var nonReleased int64
	db.Model(&models.ResourceBooking{}).
		Where("turnaround_id = ? AND booking_status <> ?", 2, constants.BookingReleased).Count(&nonReleased)
	if nonReleased != 0 {
		t.Fatalf("non-released bookings = %d, want 0", nonReleased)
	}
}

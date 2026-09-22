package constructors

import (
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

const timeLayout = "2006-01-02 15:04"

// BuildTaskBlocker 构造未 SIGNED 任务的阻塞项。
func BuildTaskBlocker(task models.GroundTask) types.TaskBlockerDTO {
	return types.TaskBlockerDTO{
		ID:          task.ID,
		TaskType:    task.TaskType,
		TeamID:      task.TeamID,
		Status:      task.Status,
		Deadline:    task.Deadline.Format(timeLayout),
		BlockerNote: task.BlockerNote,
		Reason:      fmt.Sprintf("任务 %s 当前状态 %s，放行前必须完成签收（SIGNED）", task.TaskType, task.Status),
	}
}

// BuildDelayBlocker 构造未关闭延误的阻塞项。
func BuildDelayBlocker(delay models.DelayEvent) types.DelayBlockerDTO {
	return types.DelayBlockerDTO{
		ID:                 delay.ID,
		DelayType:          delay.DelayType,
		Minutes:            delay.Minutes,
		ResponsibilityTeam: delay.ResponsibilityTeam,
		RootCause:          delay.RootCause,
		Reason:             fmt.Sprintf("延误事件 #%d（%s，+%d 分钟）尚未关闭", delay.ID, delay.DelayType, delay.Minutes),
	}
}

// BuildBookingBlocker 构造无法被放行事务安全自动释放的预约阻塞项。
// 普通 ACTIVE 会在放行通过时一次性释放，不构成阻塞；CONFLICT 必须人工先解决。
func BuildBookingBlocker(booking models.ResourceBooking) types.BookingBlockerDTO {
	resourceCode := ""
	if booking.Resource != nil {
		resourceCode = booking.Resource.ResourceCode
	}
	reason := fmt.Sprintf("资源 %s 的预约状态异常 %s，放行前必须先解决", resourceCode, booking.BookingStatus)
	if booking.BookingStatus == constants.BookingConflict {
		reason = fmt.Sprintf("资源 %s 的预约存在冲突未解决：%s", resourceCode, booking.ConflictReason)
	}
	return types.BookingBlockerDTO{
		ID:             booking.ID,
		ResourceID:     booking.ResourceID,
		ResourceCode:   resourceCode,
		TaskID:         booking.TaskID,
		BookingStatus:  booking.BookingStatus,
		ConflictReason: booking.ConflictReason,
		Reason:         reason,
	}
}

// BuildGate 聚合放行门禁快照：任务 / 延误 / 预约三类阻塞清单。
func BuildGate(t models.FlightTurnaround, tasks []models.GroundTask, delays []models.DelayEvent, bookings []models.ResourceBooking) types.ReleaseGateDTO {
	gate := types.ReleaseGateDTO{
		UnsignedTasks:    []types.TaskBlockerDTO{},
		OpenDelays:       []types.DelayBlockerDTO{},
		BlockingBookings: []types.BookingBlockerDTO{},
	}
	for _, task := range tasks {
		if task.Status != constants.TaskSigned {
			gate.UnsignedTasks = append(gate.UnsignedTasks, BuildTaskBlocker(task))
		}
	}
	for _, delay := range delays {
		if delay.ResolvedAt == nil {
			gate.OpenDelays = append(gate.OpenDelays, BuildDelayBlocker(delay))
		}
	}
	// 普通 ACTIVE 预约由放行事务一次性释放（放行后不再 ACTIVE），不阻塞；
	// 只有 CONFLICT 这类无法自动释放的预约才进入阻塞清单。
	for _, booking := range bookings {
		if booking.BookingStatus == constants.BookingConflict {
			gate.BlockingBookings = append(gate.BlockingBookings, BuildBookingBlocker(booking))
		}
	}
	gate.Releasable = len(gate.UnsignedTasks) == 0 && len(gate.OpenDelays) == 0 && len(gate.BlockingBookings) == 0
	return gate
}

// NewFlightTurnaroundResponse 构造航班过站响应（含进度与门禁明细）。
func NewFlightTurnaroundResponse(t models.FlightTurnaround) types.FlightTurnaroundDTO {
	tasks := t.Tasks
	delays := t.Delays
	bookings := t.Bookings

	progress := types.TurnaroundProgressDTO{TotalTasks: len(tasks)}
	for _, task := range tasks {
		if task.Status == constants.TaskSigned {
			progress.SignedTasks++
		}
	}
	for _, delay := range delays {
		if delay.ResolvedAt == nil {
			progress.OpenDelays++
		}
	}
	for _, booking := range bookings {
		if booking.BookingStatus == constants.BookingActive {
			progress.ActiveBookings++
		}
	}

	readyAt := (*time.Time)(nil)
	if t.ReadyAt != nil {
		readyAt = t.ReadyAt
	}

	return types.FlightTurnaroundDTO{
		ID:               t.ID,
		FlightNo:         t.FlightNo,
		AircraftReg:      t.AircraftReg,
		StandNo:          t.StandNo,
		ArrivalTime:      t.ArrivalTime,
		DepartureTime:    t.DepartureTime,
		TurnaroundStatus: t.TurnaroundStatus,
		DelayReason:      t.DelayReason,
		ReadyAt:          readyAt,
		Progress:         progress,
		Gate:             BuildGate(t, tasks, delays, bookings),
	}
}

// NewFlightTurnaroundListResponse 批量构造。
func NewFlightTurnaroundListResponse(rows []models.FlightTurnaround) []types.FlightTurnaroundDTO {
	out := make([]types.FlightTurnaroundDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewFlightTurnaroundResponse(row))
	}
	return out
}

package services

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func BuildTimeline(flight models.FlightTurnaround, tasks []models.GroundTask, bookings []models.ResourceBooking) []types.TimelineEvent {
	timeline := []types.TimelineEvent{
		{Code: "ARRIVAL", Label: "航班到达", OccurredAt: flight.ArrivalTime, Status: constants.TurnaroundArriving},
	}
	activeBookings := 0
	for _, booking := range bookings {
		if booking.BookingStatus == constants.BookingActive {
			activeBookings++
		}
	}
	timeline = append(timeline, types.TimelineEvent{Code: "RESOURCE_OCCUPANCY", Label: "资源占用", OccurredAt: flight.ArrivalTime, Status: bookingTimelineStatus(activeBookings, len(bookings))})
	for _, task := range tasks {
		occurredAt := task.PlannedStart
		if task.ActualFinish != nil {
			occurredAt = *task.ActualFinish
		}
		timeline = append(timeline, types.TimelineEvent{Code: "TASK_" + task.TaskType, Label: task.TaskType, OccurredAt: occurredAt, Status: task.Status})
	}
	if flight.TurnaroundStatus == constants.TurnaroundReady || flight.TurnaroundStatus == constants.TurnaroundDeparted {
		readyAt := flight.DepartureTime
		if flight.ReadyAt != nil {
			readyAt = *flight.ReadyAt
		}
		timeline = append(timeline, types.TimelineEvent{Code: "READY", Label: "过站放行", OccurredAt: readyAt, Status: constants.TurnaroundReady})
	}
	return timeline
}

func bookingTimelineStatus(active, total int) string {
	if total == 0 {
		return "NONE"
	}
	if active == 0 {
		return constants.BookingReleased
	}
	return constants.BookingActive
}

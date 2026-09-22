package services

import (
	"errors"
	"net/http"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

type FlightTurnaroundService struct {
	store *repositories.Store
}

func NewFlightTurnaroundService(store *repositories.Store) *FlightTurnaroundService {
	return &FlightTurnaroundService{store: store}
}

func (s *FlightTurnaroundService) List() []types.FlightTurnaroundDetail {
	snapshot := s.store.Snapshot()
	rows := make([]types.FlightTurnaroundDetail, 0, len(snapshot.Flights))
	for _, flight := range snapshot.Flights {
		tasks := filterTasks(snapshot.Tasks, flight.ID)
		bookings := filterBookings(snapshot.Bookings, flight.ID)
		rows = append(rows, constructors.NewFlightDetail(flight, BuildTimeline(flight, tasks, bookings), completionRate(tasks)))
	}
	return rows
}

func (s *FlightTurnaroundService) Get(id int) (types.FlightTurnaroundDetail, bool) {
	flight, found := s.store.GetFlight(id)
	if !found {
		return types.FlightTurnaroundDetail{}, false
	}
	snapshot := s.store.Snapshot()
	return constructors.NewFlightDetail(flight, BuildTimeline(flight, filterTasks(snapshot.Tasks, id), filterBookings(snapshot.Bookings, id)), completionRate(filterTasks(snapshot.Tasks, id))), true
}

func (s *FlightTurnaroundService) Release(id int) (types.ReleaseResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var response types.ReleaseResponse
	var serviceErr *ReleaseError

	err := s.store.Update(func(data *repositories.Data) error {
		flightIndex := -1
		for i := range data.Flights {
			if data.Flights[i].ID == id {
				flightIndex = i
				break
			}
		}
		if flightIndex < 0 {
			serviceErr = &ReleaseError{Code: constants.TurnaroundNotFound, Message: constants.TurnaroundNotFoundMessage, HTTPCode: http.StatusNotFound}
			return errors.New("flight not found")
		}
		if data.Flights[flightIndex].ReleaseInProgress {
			serviceErr = &ReleaseError{Code: constants.ConcurrentRelease, Message: constants.ConcurrentReleaseMessage, HTTPCode: http.StatusConflict}
			return errors.New("release in progress")
		}
		if data.Flights[flightIndex].TurnaroundStatus == constants.TurnaroundReady || data.Flights[flightIndex].TurnaroundStatus == constants.TurnaroundDeparted {
			serviceErr = &ReleaseError{Code: constants.TurnaroundAlreadyReleased, Message: constants.TurnaroundAlreadyReleasedMessage, HTTPCode: http.StatusConflict}
			return errors.New("already released")
		}

		data.Flights[flightIndex].ReleaseInProgress = true
		blockers := collectBlockers(data, id)
		if len(blockers.Tasks) > 0 || len(blockers.Delays) > 0 || len(blockers.Bookings) > 0 {
			data.Flights[flightIndex].ReleaseInProgress = false
			serviceErr = &ReleaseError{Code: constants.ReleaseBlocked, Message: constants.ReleaseBlockedMessage, HTTPCode: http.StatusConflict, Blockers: &blockers}
			return errors.New("release blocked")
		}

		releasedCount := 0
		releasedResources := map[int]bool{}
		for i := range data.Tasks {
			if data.Tasks[i].TurnaroundID == id {
				data.Tasks[i].TimelineSyncedAt = &now
			}
		}
		for i := range data.Bookings {
			if data.Bookings[i].TurnaroundID == id {
				releasedResources[data.Bookings[i].ResourceID] = true
				constructors.ReleaseBooking(&data.Bookings[i], now)
				releasedCount++
			}
		}
		activeByResource := map[int]int{}
		for _, booking := range data.Bookings {
			if booking.BookingStatus == constants.BookingActive {
				activeByResource[booking.ResourceID]++
			}
		}
		for i := range data.Resources {
			if releasedResources[data.Resources[i].ID] && data.Resources[i].AvailabilityStatus == constants.ResourceBooked && activeByResource[data.Resources[i].ID] == 0 {
				data.Resources[i].AvailabilityStatus = constants.ResourceAvailable
			}
		}
		data.Flights[flightIndex].TurnaroundStatus = constants.TurnaroundReady
		data.Flights[flightIndex].ReadyAt = &now
		data.Flights[flightIndex].TimelineSyncedAt = &now
		data.Flights[flightIndex].DelayReason = ""
		data.Flights[flightIndex].ReleaseInProgress = false
		data.Logs = append(data.Logs, models.AuditLog{
			ID:         len(data.Logs) + 1,
			Actor:      "dispatcher",
			Action:     "FLIGHT_TURNAROUND_RELEASE",
			TargetType: "FlightTurnaround",
			TargetID:   flightLabel(data.Flights[flightIndex]),
			CreatedAt:  now,
		})

		flight := data.Flights[flightIndex]
		tasks := filterTasks(data.Tasks, id)
		bookings := filterBookings(data.Bookings, id)
		detail := constructors.NewFlightDetail(flight, BuildTimeline(flight, tasks, bookings), completionRate(tasks))
		response = constructors.NewReleaseResponse(detail, tasks, bookings, releasedCount)
		return nil
	})
	if serviceErr != nil {
		return types.ReleaseResponse{}, serviceErr
	}
	if err != nil {
		return types.ReleaseResponse{}, err
	}
	return response, nil
}

func collectBlockers(data *repositories.Data, turnaroundID int) types.ReleaseBlockers {
	blockers := types.ReleaseBlockers{Tasks: []types.TaskBlocker{}, Delays: []types.DelayBlocker{}, Bookings: []types.BookingBlocker{}}
	resourceCodes := constructors.ResourceCodeMap(data.Resources)
	for _, task := range data.Tasks {
		if task.TurnaroundID != turnaroundID {
			continue
		}
		if task.Status != constants.GroundTaskSigned {
			blockers.Tasks = append(blockers.Tasks, constructors.NewTaskBlocker(task, "任务尚未 SIGNED"))
		}
	}
	for _, delay := range data.Delays {
		if delay.TurnaroundID == turnaroundID && delay.ResolvedAt == nil {
			blockers.Delays = append(blockers.Delays, constructors.NewDelayBlocker(delay, "延误事件尚未关闭"))
		}
	}
	for _, booking := range data.Bookings {
		if booking.TurnaroundID != turnaroundID {
			continue
		}
		if booking.BookingStatus == constants.BookingConflict {
			blockers.Bookings = append(blockers.Bookings, constructors.NewBookingBlocker(booking, resourceCodes[booking.ResourceID], "资源预约存在未处理冲突"))
		}
	}
	return blockers
}

func filterTasks(rows []models.GroundTask, turnaroundID int) []models.GroundTask {
	result := []models.GroundTask{}
	for _, row := range rows {
		if row.TurnaroundID == turnaroundID {
			result = append(result, row)
		}
	}
	return result
}

func filterBookings(rows []models.ResourceBooking, turnaroundID int) []models.ResourceBooking {
	result := []models.ResourceBooking{}
	for _, row := range rows {
		if row.TurnaroundID == turnaroundID {
			result = append(result, row)
		}
	}
	return result
}

func completionRate(tasks []models.GroundTask) int {
	if len(tasks) == 0 {
		return 100
	}
	signed := 0
	for _, task := range tasks {
		if task.Status == constants.GroundTaskSigned {
			signed++
		}
	}
	return signed * 100 / len(tasks)
}

func flightLabel(flight models.FlightTurnaround) string {
	return flight.FlightNo + "/" + flight.AircraftReg
}

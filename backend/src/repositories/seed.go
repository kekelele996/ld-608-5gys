package repositories

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
)

func SeedData() Data {
	arrivalOne := "2026-09-22T08:00:00Z"
	departureOne := "2026-09-22T09:10:00Z"
	finishedOne := "2026-09-22T09:00:00Z"
	arrivalTwo := "2026-09-22T10:00:00Z"
	departureTwo := "2026-09-22T11:10:00Z"
	arrivalThree := "2026-09-22T12:00:00Z"
	departureThree := "2026-09-22T13:10:00Z"

	return Data{
		Flights: []models.FlightTurnaround{
			{ID: 1, FlightNo: "CA1801", AircraftReg: "B-2026", StandNo: "A12", ArrivalTime: arrivalOne, DepartureTime: departureOne, TurnaroundStatus: constants.TurnaroundInService, DelayReason: ""},
			{ID: 2, FlightNo: "MU5102", AircraftReg: "B-5088", StandNo: "B07", ArrivalTime: arrivalTwo, DepartureTime: departureTwo, TurnaroundStatus: constants.TurnaroundInService, DelayReason: "行李等待"},
			{ID: 3, FlightNo: "CZ3120", AircraftReg: "B-6310", StandNo: "C03", ArrivalTime: arrivalThree, DepartureTime: departureThree, TurnaroundStatus: constants.TurnaroundOnStand, DelayReason: ""},
		},
		Tasks: []models.GroundTask{
			{ID: 101, TurnaroundID: 1, TaskType: "CLEANING", TeamID: 1, PlannedStart: "2026-09-22T08:05:00Z", Deadline: "2026-09-22T08:40:00Z", ActualFinish: &finishedOne, Status: constants.GroundTaskSigned, BlockerNote: ""},
			{ID: 102, TurnaroundID: 1, TaskType: "BAGGAGE", TeamID: 2, PlannedStart: "2026-09-22T08:05:00Z", Deadline: "2026-09-22T08:50:00Z", ActualFinish: &finishedOne, Status: constants.GroundTaskSigned, BlockerNote: ""},
			{ID: 201, TurnaroundID: 2, TaskType: "CLEANING", TeamID: 1, PlannedStart: "2026-09-22T10:05:00Z", Deadline: "2026-09-22T10:40:00Z", ActualFinish: nil, Status: constants.GroundTaskSigned, BlockerNote: ""},
			{ID: 202, TurnaroundID: 2, TaskType: "REFUEL", TeamID: 3, PlannedStart: "2026-09-22T10:10:00Z", Deadline: "2026-09-22T10:50:00Z", ActualFinish: nil, Status: constants.GroundTaskBlocked, BlockerNote: "加油栓井占用"},
			{ID: 301, TurnaroundID: 3, TaskType: "CATERING", TeamID: 4, PlannedStart: "2026-09-22T12:05:00Z", Deadline: "2026-09-22T12:40:00Z", ActualFinish: nil, Status: constants.GroundTaskSigned, BlockerNote: ""},
		},
		Resources: []models.GroundResource{
			{ID: 1, ResourceCode: "CLN-A12", ResourceType: "CLEANING", Location: "A12", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: "2026-10-01T00:00:00Z", OwnerTeam: "清洁一组"},
			{ID: 2, ResourceCode: "BAG-07", ResourceType: "BAGGAGE", Location: "B07", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: "2026-10-02T00:00:00Z", OwnerTeam: "行李班组"},
			{ID: 3, ResourceCode: "FUEL-B07", ResourceType: "REFUEL", Location: "B07", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: "2026-10-03T00:00:00Z", OwnerTeam: "加油班组"},
			{ID: 4, ResourceCode: "CAT-C03", ResourceType: "CATERING", Location: "C03", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: "2026-10-04T00:00:00Z", OwnerTeam: "配餐班组"},
		},
		Bookings: []models.ResourceBooking{
			{ID: 1001, ResourceID: 1, TurnaroundID: 1, TaskID: 101, StartTime: "2026-09-22T08:05:00Z", EndTime: "2026-09-22T08:40:00Z", BookingStatus: constants.BookingActive, ConflictReason: ""},
			{ID: 1002, ResourceID: 2, TurnaroundID: 1, TaskID: 102, StartTime: "2026-09-22T08:05:00Z", EndTime: "2026-09-22T08:50:00Z", BookingStatus: constants.BookingActive, ConflictReason: ""},
			{ID: 2001, ResourceID: 3, TurnaroundID: 2, TaskID: 202, StartTime: "2026-09-22T10:10:00Z", EndTime: "2026-09-22T10:50:00Z", BookingStatus: constants.BookingConflict, ConflictReason: "加油栓井与邻机位冲突"},
			{ID: 3001, ResourceID: 4, TurnaroundID: 3, TaskID: 301, StartTime: "2026-09-22T12:05:00Z", EndTime: "2026-09-22T12:40:00Z", BookingStatus: constants.BookingActive, ConflictReason: ""},
		},
		Delays: []models.DelayEvent{
			{ID: 1, TurnaroundID: 2, DelayType: "GROUND_HANDLING", Minutes: 25, RootCause: "加油资源冲突", ResponsibilityTeam: "加油班组", ResolvedAt: nil},
		},
		Logs: []models.AuditLog{},
	}
}

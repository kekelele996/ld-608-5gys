package config

import (
	"log"
	"time"

	"gorm.io/gorm"

	"groundTurn/src/constants"
	"groundTurn/src/models"
)

const timeFmt = "2006-01-02 15:04:05"

func mustTime(value string) time.Time {
	t, err := time.ParseInLocation(timeFmt, value, time.Local)
	if err != nil {
		log.Fatalf("bad seed time %q: %v", value, err)
	}
	return t
}

// SeedIfEmpty 空库时灌入演示数据，覆盖放行联动三类典型场景。
func SeedIfEmpty(db *gorm.DB) {
	var count int64
	db.Model(&models.FlightTurnaround{}).Count(&count)
	if count > 0 {
		return
	}

	base := mustTime("2026-09-22 08:00:00")

	// 场景 1：MU518 门禁全绿 —— 任务全 SIGNED、无未关闭延误、预约均 ACTIVE（放行时一次性释放）。
	ta1 := models.FlightTurnaround{
		FlightNo: "MU518", AircraftReg: "B-2201", StandNo: "103",
		ArrivalTime: base, DepartureTime: base.Add(55 * time.Minute),
		TurnaroundStatus: constants.TurnaroundInService, DelayReason: "",
	}
	// 场景 2：CA183 三类阻塞齐全 —— 有未签收任务、未关闭延误、CONFLICT 预约。
	ta2 := models.FlightTurnaround{
		FlightNo: "CA183", AircraftReg: "B-5520", StandNo: "117",
		ArrivalTime: base.Add(40 * time.Minute), DepartureTime: base.Add(100 * time.Minute),
		TurnaroundStatus: constants.TurnaroundInService, DelayReason: "货舱装卸等待",
	}
	// 场景 3：CZ302 仅有延误阻塞 —— 任务已全 SIGNED，预约 ACTIVE，但挂一条未关闭延误。
	ta3 := models.FlightTurnaround{
		FlightNo: "CZ302", AircraftReg: "B-8819", StandNo: "109",
		ArrivalTime: base.Add(70 * time.Minute), DepartureTime: base.Add(135 * time.Minute),
		TurnaroundStatus: constants.TurnaroundInService, DelayReason: "航油车排队",
	}
	if err := db.Create([]*models.FlightTurnaround{&ta1, &ta2, &ta3}).Error; err != nil {
		log.Fatalf("seed turnarounds: %v", err)
	}

	resources := []models.GroundResource{
		{ResourceCode: "GPU-01", ResourceType: "GPU", Location: "T2-103", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: base.AddDate(0, 1, 0), OwnerTeam: "机务一组"},
		{ResourceCode: "CART-CAT", ResourceType: "CATERING", Location: "T2-117", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: base.AddDate(0, 1, 0), OwnerTeam: "配餐组"},
		{ResourceCode: "FUEL-07", ResourceType: "REFUEL", Location: "T2-109", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: base.AddDate(0, 1, 0), OwnerTeam: "加油组"},
		{ResourceCode: "BELT-12", ResourceType: "BAGGAGE", Location: "T2-117", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: base.AddDate(0, 1, 0), OwnerTeam: "行李组"},
		{ResourceCode: "WATER-03", ResourceType: "WATER", Location: "T2-103", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: base.AddDate(0, 1, 0), OwnerTeam: "清水组"},
		{ResourceCode: "PB-02", ResourceType: "PUSHBACK", Location: "T2-103", AvailabilityStatus: constants.ResourceBooked, MaintenanceDueAt: base.AddDate(0, 1, 0), OwnerTeam: "牵引车组"},
	}
	if err := db.Create(&resources).Error; err != nil {
		log.Fatalf("seed resources: %v", err)
	}

	finish := base.Add(35 * time.Minute)
	tasks := []models.GroundTask{
		// TA1 全部 SIGNED
		{TurnaroundID: ta1.ID, TaskType: "CLEANING", TeamID: 1, PlannedStart: base.Add(5 * time.Minute), Deadline: base.Add(30 * time.Minute), ActualFinish: &finish, Status: constants.TaskSigned},
		{TurnaroundID: ta1.ID, TaskType: "CATERING", TeamID: 2, PlannedStart: base.Add(8 * time.Minute), Deadline: base.Add(32 * time.Minute), ActualFinish: &finish, Status: constants.TaskSigned},
		{TurnaroundID: ta1.ID, TaskType: "REFUEL", TeamID: 3, PlannedStart: base.Add(10 * time.Minute), Deadline: base.Add(40 * time.Minute), ActualFinish: &finish, Status: constants.TaskSigned},
		// TA2 三个未 SIGNED
		{TurnaroundID: ta2.ID, TaskType: "BAGGAGE", TeamID: 4, PlannedStart: base.Add(45 * time.Minute), Deadline: base.Add(70 * time.Minute), Status: constants.TaskInProgress, BlockerNote: "传送带车未到位"},
		{TurnaroundID: ta2.ID, TaskType: "WATER_SERVICE", TeamID: 5, PlannedStart: base.Add(50 * time.Minute), Deadline: base.Add(75 * time.Minute), Status: constants.TaskBlocked, BlockerNote: "清水车故障"},
		{TurnaroundID: ta2.ID, TaskType: "PUSHBACK", TeamID: 6, PlannedStart: base.Add(80 * time.Minute), Deadline: base.Add(95 * time.Minute), Status: constants.TaskDispatched},
		// TA3 全部 SIGNED
		{TurnaroundID: ta3.ID, TaskType: "BAGGAGE", TeamID: 4, PlannedStart: base.Add(75 * time.Minute), Deadline: base.Add(100 * time.Minute), ActualFinish: &finish, Status: constants.TaskSigned},
		{TurnaroundID: ta3.ID, TaskType: "REFUEL", TeamID: 3, PlannedStart: base.Add(80 * time.Minute), Deadline: base.Add(110 * time.Minute), ActualFinish: &finish, Status: constants.TaskSigned},
	}
	if err := db.Create(&tasks).Error; err != nil {
		log.Fatalf("seed tasks: %v", err)
	}

	bookings := []models.ResourceBooking{
		{ResourceID: resources[0].ID, TurnaroundID: ta1.ID, TaskID: tasks[0].ID, StartTime: base.Add(5 * time.Minute), EndTime: base.Add(40 * time.Minute), BookingStatus: constants.BookingActive},
		{ResourceID: resources[4].ID, TurnaroundID: ta1.ID, TaskID: tasks[1].ID, StartTime: base.Add(8 * time.Minute), EndTime: base.Add(42 * time.Minute), BookingStatus: constants.BookingActive},
		{ResourceID: resources[5].ID, TurnaroundID: ta1.ID, TaskID: tasks[2].ID, StartTime: base.Add(10 * time.Minute), EndTime: base.Add(48 * time.Minute), BookingStatus: constants.BookingActive},
		// TA2：一条 ACTIVE + 一条 CONFLICT
		{ResourceID: resources[1].ID, TurnaroundID: ta2.ID, TaskID: tasks[3].ID, StartTime: base.Add(45 * time.Minute), EndTime: base.Add(80 * time.Minute), BookingStatus: constants.BookingActive},
		{ResourceID: resources[3].ID, TurnaroundID: ta2.ID, TaskID: tasks[4].ID, StartTime: base.Add(50 * time.Minute), EndTime: base.Add(85 * time.Minute), BookingStatus: constants.BookingConflict, ConflictReason: "与 CA166 的行李传送带占用重叠"},
		// TA3：ACTIVE
		{ResourceID: resources[2].ID, TurnaroundID: ta3.ID, TaskID: tasks[7].ID, StartTime: base.Add(80 * time.Minute), EndTime: base.Add(115 * time.Minute), BookingStatus: constants.BookingActive},
	}
	if err := db.Create(&bookings).Error; err != nil {
		log.Fatalf("seed bookings: %v", err)
	}

	delays := []models.DelayEvent{
		{TurnaroundID: ta2.ID, DelayType: "BAGGAGE", Minutes: 18, RootCause: "行李分拣系统卡包", ResponsibilityTeam: "行李组"},
		{TurnaroundID: ta3.ID, DelayType: "REFUEL", Minutes: 12, RootCause: "航油车调度排队", ResponsibilityTeam: "加油组"},
		// TA1 一条已关闭延误，不计阻塞。
		{TurnaroundID: ta1.ID, DelayType: "CATERING", Minutes: 5, RootCause: "配餐晚到，已追回", ResponsibilityTeam: "配餐组", ResolvedAt: &finish},
	}
	if err := db.Create(&delays).Error; err != nil {
		log.Fatalf("seed delays: %v", err)
	}

	log.Printf("seed complete: %d turnarounds, %d tasks, %d bookings, %d delays",
		3, len(tasks), len(bookings), len(delays))
}

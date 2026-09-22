package models

import "time"

// FlightTurnaround 航班过站。贯穿 models/repositories/services/controllers/routes/前端全链路。
type FlightTurnaround struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	FlightNo         string     `gorm:"column:flight_no;size:32" json:"flight_no"`
	AircraftReg      string     `gorm:"column:aircraft_reg;size:32" json:"aircraft_reg"`
	StandNo          string     `gorm:"column:stand_no;size:32" json:"stand_no"`
	ArrivalTime      time.Time  `gorm:"column:arrival_time" json:"arrival_time"`
	DepartureTime    time.Time  `gorm:"column:departure_time" json:"departure_time"`
	TurnaroundStatus string     `gorm:"column:turnaround_status;size:32;index" json:"turnaround_status"`
	DelayReason      string     `gorm:"column:delay_reason;size:255" json:"delay_reason"`
	ReadyAt          *time.Time `gorm:"column:ready_at;index" json:"ready_at"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// 放行门禁聚合（非数据库列），看板直接展示阻塞明细。
	Tasks    []GroundTask      `gorm:"foreignKey:TurnaroundID" json:"tasks,omitempty"`
	Bookings []ResourceBooking `gorm:"foreignKey:TurnaroundID" json:"bookings,omitempty"`
	Delays   []DelayEvent      `gorm:"foreignKey:TurnaroundID" json:"delays,omitempty"`
}

func (FlightTurnaround) TableName() string { return "flight_turnaround" }

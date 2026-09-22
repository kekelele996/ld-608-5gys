package models

import "time"

// GroundResource 保障资源台账。
type GroundResource struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceCode       string    `gorm:"column:resource_code;size:32" json:"resource_code"`
	ResourceType       string    `gorm:"column:resource_type;size:32" json:"resource_type"`
	Location           string    `gorm:"column:location;size:64" json:"location"`
	AvailabilityStatus string    `gorm:"column:availability_status;size:32;index" json:"availability_status"`
	MaintenanceDueAt   time.Time `gorm:"column:maintenance_due_at" json:"maintenance_due_at"`
	OwnerTeam          string    `gorm:"column:owner_team;size:64" json:"owner_team"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (GroundResource) TableName() string { return "ground_resource" }

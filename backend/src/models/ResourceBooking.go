package models

import "time"

// ResourceBooking 资源预约。ACTIVE 在放行事务中一次性释放；CONFLICT 是硬阻塞。
type ResourceBooking struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceID     int64      `gorm:"column:resource_id;index" json:"resource_id"`
	TurnaroundID   int64      `gorm:"column:turnaround_id;index" json:"turnaround_id"`
	TaskID         int64      `gorm:"column:task_id" json:"task_id"`
	StartTime      time.Time  `gorm:"column:start_time" json:"start_time"`
	EndTime        time.Time  `gorm:"column:end_time" json:"end_time"`
	BookingStatus  string     `gorm:"column:booking_status;size:32;index" json:"booking_status"`
	ConflictReason string     `gorm:"column:conflict_reason;size:255" json:"conflict_reason"`
	ReleasedAt     *time.Time `gorm:"column:released_at" json:"released_at"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`

	Resource *GroundResource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}

func (ResourceBooking) TableName() string { return "resource_booking" }

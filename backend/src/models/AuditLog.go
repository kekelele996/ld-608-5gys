package models

import "time"

// AuditLog 操作日志：所有写操作（派工、预约、延误归因、放行）必须记录。
type AuditLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Actor      string    `gorm:"column:actor;size:64" json:"actor"`
	Action     string    `gorm:"column:action;size:64" json:"action"`
	TargetType string    `gorm:"column:target_type;size:64" json:"target_type"`
	TargetID   string    `gorm:"column:target_id;size:64" json:"target_id"`
	Detail     string    `gorm:"column:detail;size:512" json:"detail"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_log" }

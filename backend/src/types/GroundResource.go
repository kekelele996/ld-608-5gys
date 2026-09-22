package types

import "time"

type GroundResourceDTO struct {
	ID                 int64     `json:"id"`
	ResourceCode       string    `json:"resource_code"`
	ResourceType       string    `json:"resource_type"`
	Location           string    `json:"location"`
	AvailabilityStatus string    `json:"availability_status"`
	MaintenanceDueAt   time.Time `json:"maintenance_due_at"`
	OwnerTeam          string    `json:"owner_team"`
}

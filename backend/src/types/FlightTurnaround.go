package types

import "groundTurn/src/models"

type TimelineEvent struct {
	Code       string `json:"code"`
	Label      string `json:"label"`
	OccurredAt string `json:"occurred_at"`
	Status     string `json:"status"`
}

type FlightTurnaroundDetail struct {
	models.FlightTurnaround
	TaskCompletionRate int             `json:"task_completion_rate"`
	Timeline           []TimelineEvent `json:"timeline"`
}

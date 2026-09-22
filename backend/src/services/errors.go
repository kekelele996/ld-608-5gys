package services

import "groundTurn/src/types"

type ReleaseError struct {
	Code     string
	Message  string
	HTTPCode int
	Blockers *types.ReleaseBlockers
}

func (e *ReleaseError) Error() string { return e.Message }

package constants

const (
	TurnaroundArriving  = "ARRIVING"
	TurnaroundOnStand   = "ON_STAND"
	TurnaroundInService = "IN_SERVICE"
	TurnaroundReady     = "READY"
	TurnaroundDeparted  = "DEPARTED"
	TurnaroundDelayed   = "DELAYED"
)

var TurnaroundStatuses = []string{
	TurnaroundArriving,
	TurnaroundOnStand,
	TurnaroundInService,
	TurnaroundReady,
	TurnaroundDeparted,
	TurnaroundDelayed,
}

package constants

const (
	ResourceAvailable   = "AVAILABLE"
	ResourceBooked      = "BOOKED"
	ResourceMaintenance = "MAINTENANCE"
	ResourceOffline     = "OFFLINE"
)

var ResourceStatuses = []string{ResourceAvailable, ResourceBooked, ResourceMaintenance, ResourceOffline}

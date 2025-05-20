package event_attendance

type permission struct {
	Key string
	// Index all event attendances
	Index string
	// View event attendance details
	Show string
}

var Permissions = permission{
	Key:   "event-attendance",
	Index: "index",
	Show:  "show",
}

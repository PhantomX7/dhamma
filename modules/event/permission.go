package event

type permission struct {
	Key string
	// Index all events
	Index string
	// View event details
	Show string
	// Create a new event
	Create string
	// Update event information
	Update string
	// Attend an event
	Attend string
}

var Permissions = permission{
	Key:    "event",
	Index:  "index",
	Show:   "show",
	Create: "create",
	Update: "update",
	Attend: "attend",
}

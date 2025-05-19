package event

type permission struct {
	Key string
	// Index all followers
	Index string
	// View follower details
	Show string
	// Create a new follower
	Create string
	// Update follower information
	Update string
}

var Permissions = permission{
	Key:    "event",
	Index:  "index",
	Show:   "show",
	Create: "create",
	Update: "update",
}

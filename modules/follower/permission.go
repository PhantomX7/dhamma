package follower

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
	// Add a card to a follower
	AddCard string
	// Delete a card from a follower
	DeleteCard string
}

var Permissions = permission{
	Key:        "follower",
	Index:      "index",
	Show:       "show",
	Create:     "create",
	Update:     "update",
	AddCard:    "add-card",
	DeleteCard: "delete-card",
}

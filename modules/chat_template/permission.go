package chat_template

type permission struct {
	Key string
	// Index all chat templates
	Index string
	// View chat template details
	Show string
	// Create a new chat template
	Create string
	// Update chat template information
	Update string
	// Delete a chat template
	Delete string
	// Set chat template as default
	SetAsDefault string
	// Get default chat template
	GetDefault string
}

var Permissions = permission{
	Key:          "chat-template",
	Index:        "index",
	Show:         "show",
	Create:       "create",
	Update:       "update",
	Delete:       "delete",
	SetAsDefault: "set-as-default",
	GetDefault:   "get-default",
}

package auth

type permission struct {
	Key string
	// Get the current user information
	GetMe string
	// Update the current user's password
	UpdatePassword string
}

var Permissions = permission{
	Key:            "auth",
	GetMe:          "me",
	UpdatePassword: "update-password",
}

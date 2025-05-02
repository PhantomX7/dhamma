package auth

type permission struct {
	Key string
	// Update the current user's password
	UpdatePassword string
}

var Permissions = permission{
	Key:            "auth",
	UpdatePassword: "update-password",
}

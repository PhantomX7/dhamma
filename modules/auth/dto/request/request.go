package request

type SignInRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type SignUpRequest struct {
	Username string `form:"username" json:"username" binding:"required,unique=users.username"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" form:"current_password" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
}

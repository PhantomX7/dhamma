package request

type SignInRequest struct {
	Username   string  `form:"username" json:"username" binding:"required"`
	Password   string  `form:"password" json:"password" binding:"required"`
	DomainCode *string `form:"domain_code" json:"domain_code,omitempty"`
}

type SignUpRequest struct {
	Username string `form:"username" json:"username" binding:"required,unique=users.username"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" form:"current_password" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

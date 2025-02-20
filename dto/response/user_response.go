package response

import "github.com/PhantomX7/dhamma/entity"

// Response
type (
	UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	}

	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		PaginationResponse
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User `json:"users"`
		PaginationResponse
	}

	UserUpdateResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	}

	UserLoginResponse struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}

	// VerifyEmailResponse struct {
	// 	Email      string `json:"email"`
	// 	IsVerified bool   `json:"is_verified"`
	// }
)

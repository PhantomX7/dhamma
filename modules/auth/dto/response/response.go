package response

import "github.com/PhantomX7/dhamma/entity"

type AuthResponse struct {
	Token string `json:"token"`
}

type MeResponse struct {
	entity.User
}

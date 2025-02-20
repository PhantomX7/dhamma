package response

import "github.com/PhantomX7/dhamma/model"

type AuthResponse struct {
	Token string `json:"token"`
}

type MeResponse struct {
	model.User
}

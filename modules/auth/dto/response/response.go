package response

import "github.com/PhantomX7/dhamma/entity"

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MeResponse struct {
	entity.User
	DomainID uint64 `json:"domain_id"`
}

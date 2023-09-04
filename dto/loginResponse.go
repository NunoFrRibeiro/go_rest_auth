package dto

type LoginResponse struct {
	AccessToken  string `json:"access_toke"`
	RefreshToken string `json:"refresh_token"`
}

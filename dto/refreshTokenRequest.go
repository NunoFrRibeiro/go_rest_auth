package dto

import (
	"errors"

	"github.com/NunoFrRibeiro/go_rest_auth/domain"
	"github.com/golang-jwt/jwt"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (r RefreshTokenRequest) IsAccessTokenValid() *jwt.ValidationError {
	_, err := jwt.Parse(r.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, vErr) {
			return vErr
		}
	}
	return nil
}

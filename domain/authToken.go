package domain

import (
	errs "github.com/NunoFrRibeiro/go_rest_auth/err"
	"github.com/NunoFrRibeiro/go_rest_auth/logger"
	"github.com/golang-jwt/jwt"
)

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("failed to sign access token: %s", err.Error())
		return "", errs.UnexpectedError("cannont generate access token")
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, *errs.AppError) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshClaims := c.RefreshAccessTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("failed to sign refresh token: %s", err.Error())
	}
	return signedString, nil
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}

func NewAccesTokenFromRefresh(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		return "", errs.AuthenticationError("invalid or expired token")
	}
	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenCLaims()
	authToken := NewAuthToken(accessTokenClaims)

	return authToken.NewAccessToken()
}

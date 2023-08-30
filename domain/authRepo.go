package domain

import (
	errs "github.com/NunoFrRibeiro/go_rest_auth/err"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username, password string) (*Login, error)
	GenerateAndStoreRefreshTokenToStore(authToken AuthToken) (string, errs.AppError)
	VerifyRefreshToken(refreshToken string) *errs.AppError
}

type AuthRepoDB struct {
	Client *sqlx.DB
}

func NewAuthRepo(client *sqlx.DB) AuthRepoDB {
	return AuthRepoDB{
		Client: client,
	}
}

func (db AuthRepoDB) FindBy(username, password string)

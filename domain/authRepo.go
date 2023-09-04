package domain

import (
	"database/sql"

	errs "github.com/NunoFrRibeiro/go_rest_auth/err"
	"github.com/NunoFrRibeiro/go_rest_auth/logger"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username, password string) (*Login, *errs.AppError)
	GenerateAndStoreRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError)
	VerifyRefreshToken(refreshToken string) *errs.AppError
}

type AuthRepoDB struct {
	Client *sqlx.DB
}

func NewAuthRepo(client *sqlx.DB) AuthRepoDB {
	return AuthRepoDB{client}
}

func (db AuthRepoDB) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlQuery := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
				LEFT JOIN accounts a 
				ON a.customer_id = u.customer_id WHERE username = ?  and password = ?
				GROUP BY a.customer_id`
	err := db.Client.Get(&login, sqlQuery, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.AuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: %s\n", err)
			return nil, errs.UnexpectedError("unexpect database error")
		}
	}
	return &login, nil
}

func (db AuthRepoDB) VerifyRefreshToken(refreshTken string) *errs.AppError {
	sqlQuery := `SELECT refresh_token from refresh_token_store where refresh_token = ?`
	var token string
	err := db.Client.Get(&token, sqlQuery, refreshTken)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("unexpected database error: %s\n", err)
			return errs.UnexpectedError("unexpected database error")
		}
	}
	return nil
}

func (db AuthRepoDB) GenerateAndStoreRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	var (
		appError     *errs.AppError
		refreshToken string
	)

	if refreshToken, appError = authToken.newRefreshToken(); appError != nil {
		return "", appError
	}

	sqlQuery := `INSERT into refresh_token_store (refresh_token) values (?)`
	_, err := db.Client.Exec(refreshToken, sqlQuery)
	if err != nil {
		logger.Error("unexpected database error: %s\n", err)
		return "", errs.UnexpectedError("unexpect database error")
	}
	return refreshToken, nil
}

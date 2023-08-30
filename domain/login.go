package domain

import "database/sql"

type Login struct {
	Username   string         `db:"username"`
	CustomerId sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"accounts"`
	role       string         `db:"role"`
}

func (l Login) ClaimsForAccessToken() 
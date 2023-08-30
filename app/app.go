package app

import (
	"fmt"
	"os"
	"time"

	"github.com/NunoFrRibeiro/go_rest_auth/domain"
	"github.com/NunoFrRibeiro/go_rest_auth/logger"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Start() {
	sanityCheck()

	router := mux.NewRouter()
	authRepo := domain.NewAuthRepo(getClientDB())

}

func sanityCheck() {
	envVar := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_ADDRESS",
		"DB_PORT",
		"DB_NAME",
	}

	for _, i := range envVar {
		if os.Getenv(i) == "" {
			logger.Info("envvironment variable %s not defined", i)
		}
	}
}

func getClientDB() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}

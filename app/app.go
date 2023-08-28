package app

import (
	"os"

	"github.com/NunoFrRibeiro/go_rest_auth/logger"
)

func Start() {
	sanityCheck()
}

func sanityCheck() {
	envVar := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_ADDRESS",
		"BD_PORT",
		"DB_NAME",
	}

	for _, i := range envVar {
		if os.Getenv(i) == "" {
			logger.Info("envvironment variable %s not defined", i)
		}
	}
}

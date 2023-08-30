package main

import (
	"github.com/NunoFrRibeiro/go_rest_auth/app"
	"github.com/NunoFrRibeiro/go_rest_auth/logger"
)

func main() {
	logger.Info("Starting Server...")
	app.Start()
}

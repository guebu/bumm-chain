package main

import (
	"github.com/guebu/common-utils/logger"
	"go.mod/app"
)

func main() {
	logger.Info("Start of Application stock-api", "App:stock-api", "Layer:app", "Status:Open")
	app.StartApplication()
	logger.Info("End of Application stock-api", "App:stock-api", "Layer:app", "Status:End")
}





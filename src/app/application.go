package app

import(
	"fmt"
	"github.com/guebu/common-utils/logger"
)

func StartApplication() {
	logger.Info("Starting application...", "Layer:app", "Func:StartApplication", "Status:Start")
	fmt.Println("Hallo!!")
	logger.Info("Starting application...", "Layer:app", "Func:StartApplication", "Status:End")
}

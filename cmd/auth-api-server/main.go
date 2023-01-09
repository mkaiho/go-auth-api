package main

import (
	"github.com/mkaiho/go-auth-api/util"
)

func init() {
	util.InitGLogger(
		util.OptionLoggerLevel(util.LoggerLevelDebug),
		util.OptionLoggerFormat(util.LoggerFormatJSON),
	)
}

func main() {
	logger := util.GLogger()

	logger.Info("completed")
}

package main

import (
	"github.com/HasanShahjahan/sidecar-service/internal/config"
	"github.com/HasanShahjahan/sidecar-service/internal/logger"
	"github.com/joho/godotenv"
)

const (
	logTag = "Start"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logging.Error(logTag, "Error getting env, not coming through %+v", err)
	} else {
		logging.Info(logTag, "We are getting the env values")
	}
	if err := config.LoadJSONConfig(config.Config); err != nil {
		logging.Fatal(logTag, "unable to load configuration. error=%+v", err)
	}
	logging.SetLogLevel(config.Config.LogLevel)
}



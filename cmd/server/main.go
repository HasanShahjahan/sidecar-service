package main

import (
	"github.com/HasanShahjahan/sidecar-service/internal/config"
	"github.com/HasanShahjahan/sidecar-service/internal/logger"
	"github.com/joho/godotenv"
)

const (
	logTag = "Start"
)

var allowedList []string

func main() {
	allowedList = []string{
		"/company/",
		"/company/{id}",
		"/company/account",
		"/account",
		"/account/{id}",
		"/{id}",
		"/account/{id}/user",
		"/tenant/account/blocked",
	}
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

func ValidatePath(path string) bool {
	// Some magic here...
	return false
}

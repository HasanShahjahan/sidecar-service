package main

import (
	"github.com/HasanShahjahan/sidecar-service/internal/config"
	"github.com/HasanShahjahan/sidecar-service/internal/logger"
	"github.com/HasanShahjahan/sidecar-service/internal/proxy"
	"github.com/joho/godotenv"
	"log"
	"regexp"
	"strings"
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

	var allowedList = []*regexp.Regexp{
		regexp.MustCompile(`^company$`),
		regexp.MustCompile(strings.Replace(`^company/{id}$`, "{id}", proxy.IdFormat, 1)),
		regexp.MustCompile(`^company/account$`),
		regexp.MustCompile(`^account$`),
		regexp.MustCompile(strings.Replace(`^account/{id}$`, "{id}", proxy.IdFormat, 1)),
		regexp.MustCompile(strings.Replace(`^{id}$`, "{id}", proxy.IdFormat, 1)),
		regexp.MustCompile(strings.Replace(`^account/{id}/user$`, "{id}", proxy.IdFormat, 1)),
		regexp.MustCompile(`^tenant/account/blocked$`),
	}

	proxy := proxy.NewProxy(allowedList)
	isValid :=proxy.ValidatePath("company")
	log.Println(isValid)
}



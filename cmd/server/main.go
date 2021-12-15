package main

import (
	"fmt"
	"github.com/HasanShahjahan/sidecar-service/internal/config"
	"log"

	"github.com/joho/godotenv"
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
	fmt.Println("Hasan")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error")
	} else {
		fmt.Println("We are getting the env values")
	}

	if err := config.LoadJSONConfig(config.Config); err != nil {
		log.Fatal("Json load error")
	}

	fmt.Println(config.Config.LogLevel)

}

func ValidatePath(path string) bool {
	// Some magic here...
	return false
}

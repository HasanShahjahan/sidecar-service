package main

import "fmt"

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
}

func ValidatePath(path string) bool {
	// Some magic here...
	return false
}

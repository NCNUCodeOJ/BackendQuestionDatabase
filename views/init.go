package views

import "os"

var needLog = false
var userHost = os.Getenv("USER_HOST")
var privateBaseURL = "/api/private/v1"

// Setup The api init function is called when the application starts.
func Setup() {
	if os.Getenv("LOG") == "1" {
		needLog = true
	}

}

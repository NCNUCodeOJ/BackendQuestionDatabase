package views

import "os"

var needLog = false

// Setup The api init function is called when the application starts.
func Setup() {
	if os.Getenv("LOG") == "1" {
		needLog = true
	}
}

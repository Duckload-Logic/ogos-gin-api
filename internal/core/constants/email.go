package constants

import (
	"fmt"
	"os"
)

func FromEmail() string {
	sender := "noreply"
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "localhost"
	}
	return fmt.Sprintf("%s@%s", sender, domain)
}

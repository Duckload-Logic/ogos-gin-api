package builders

import (
	"os"
)

func BuildFileURL(relativePaths ...string) string {
	basePath := os.Getenv("UPLOAD_DIR")
	result := basePath
	for _, path := range relativePaths {
		result = result + "/" + path
	}

	return result
}

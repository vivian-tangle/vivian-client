package tools

import (
	"os"
)

// DBExists check if the badger DB exists
func DBExists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

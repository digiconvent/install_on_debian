package utils

import "os"

func FileExists(uri string) bool {
	_, err := os.Stat(uri)

	return err == nil
}

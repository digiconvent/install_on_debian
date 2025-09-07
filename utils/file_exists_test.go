package utils_test

import (
	"os"
	"testing"

	"github.com/digiconvent/install_on_debian/utils"
)

func TestFileExists(t *testing.T) {
	existingFile, _ := os.Executable()
	nonExistingFile := existingFile + "made-up-shit"

	if utils.FileExists(existingFile) == false {
		t.Fatal("expected", existingFile, "to exist")
	}

	if utils.FileExists(nonExistingFile) {
		t.Fatal("expected", nonExistingFile, "not to exist")
	}
}

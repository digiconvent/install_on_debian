package binary

import (
	"path"

	"github.com/digiconvent/install_on_debian/utils"
)

func TargetBinaryExists(name string) bool {
	return utils.FileExists(path.Join("/home", name, "main"))
}

package user

import (
	"os/user"
	"path"

	"github.com/digiconvent/install_on_debian/utils"
)

func UserExists(name string) bool {
	_, err := user.Lookup(name)
	return err == nil
}

func HomeFolderExists(name string) bool {
	return utils.FileExists(path.Join("/home", name))
}

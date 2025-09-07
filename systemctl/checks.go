package systemctl

import (
	"path"

	"github.com/digiconvent/install_on_debian/utils"
)

func ServiceFileExists(name string) bool {
	return utils.FileExists(path.Join("/etc", "systemd", "system", name+".service"))
}

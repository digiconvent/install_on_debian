package install_on_debian

import (
	"github.com/digiconvent/install_on_debian/binary"
	"github.com/digiconvent/install_on_debian/systemctl"
	user "github.com/digiconvent/install_on_debian/user"
)

func InstallThisBinary(name string) (systemctl.StartedService, error) {
	u, err := user.CreateOrGetUser(name)
	if err != nil {
		return nil, err
	}

	sysCtl, err := systemctl.Get(name)
	if err != nil {
		return nil, err
	}

	bin := binary.New(name)
	err = bin.HardLinkToHome()
	if err != nil {
		return nil, err
	}

	sysCtl.User = u
	var startedService systemctl.StartedService
	if !sysCtl.IsInstalled() {
		idleService, err := sysCtl.Install("")
		if err != nil {
			return nil, err
		}
		startedService, err = idleService.Start()
		if err != nil {
			return nil, err
		}
	} else {
		startedService, err = sysCtl.Start()
		if err != nil {
			return nil, err
		}
	}

	return startedService, nil
}

func UninstallThisBinary(name string) error {
	sysCtl, err := systemctl.Get(name)
	if err != nil {
		return err
	}

	if sysCtl.IsRunning() {
		_, err := sysCtl.Stop()
		if err != nil {
			return err
		}
	}

	_, err = sysCtl.Uninstall()
	if err != nil {
		return err
	}

	sysCtl.User.Delete()

	return nil
}

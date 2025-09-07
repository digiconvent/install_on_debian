package install_on_debian

import (
	"github.com/digiconvent/install_on_debian/binary"
	"github.com/digiconvent/install_on_debian/systemctl"
	user "github.com/digiconvent/install_on_debian/user"
)

type Binary struct {
	name string
}

func NewBinary(name string) *Binary {
	return &Binary{
		name: name,
	}
}

func (b *Binary) Install() (systemctl.StartedService, error) {
	u, err := user.CreateOrGetUser(b.name)
	if err != nil {
		return nil, err
	}

	sysCtl, err := systemctl.Get(b.name)
	if err != nil {
		return nil, err
	}

	bin := binary.New(b.name)
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

func (b *Binary) Uninstall() error {
	sysCtl, err := systemctl.Get(b.name)
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

func (b *Binary) IsInstalled() (bool, error) {
	c, err := systemctl.Get(b.name)
	if err != nil {
		return false, err
	}
	return c.IsInstalled(), nil
}

func (b *Binary) IsRunning() (bool, error) {
	c, err := systemctl.Get(b.name)

	if err != nil {
		return false, err
	}
	return c.IsRunning(), nil
}

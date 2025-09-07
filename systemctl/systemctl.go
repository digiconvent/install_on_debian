package systemctl

import (
	"errors"
	"os"
	"os/exec"

	"github.com/digiconvent/install_on_debian/user"
	"github.com/digiconvent/install_on_debian/utils"
)

type SystemCtlI interface {
	Install(unitFile string) (IdleService, error)
	Uninstall() (UninstalledService, error)
	Start() (StartedService, error)
	Stop() (IdleService, error)
	DeleteAccount() error

	IsRunning() bool
	IsInstalled() bool
	refreshStatus() (*ServiceStatus, error)

	reload() error
}

type UninstalledService interface {
	Install(unitFile string) (IdleService, error)
	DeleteAccount() error
}

type StartedService interface {
	Stop() (IdleService, error)
}
type IdleService interface {
	Uninstall() (UninstalledService, error)
	Start() (StartedService, error)
}

type SystemCtl struct {
	serviceName string
	User        *user.OsUserAccount
	status      *ServiceStatus
}

func (s *SystemCtl) DeleteAccount() error {
	return s.User.Delete()
}

func (s *SystemCtl) reload() error {
	result, err := utils.Execute("systemctl daemon-reload")
	if err != nil {
		return errors.New(result + err.Error())
	}
	return nil
}

func (s *SystemCtl) Uninstall() (UninstalledService, error) {
	if _, err := os.Stat(servicePath(s.serviceName)); err == nil {
		err := os.Remove(servicePath(s.serviceName))
		if err != nil {
			return nil, err
		}
	}

	if s.User != nil {
		err := s.User.Delete()
		if err != nil {
			return nil, err
		}
	}

	s.refreshStatus()
	return s, nil
}

func (s *SystemCtl) Install(unitFile string) (IdleService, error) {
	var err error
	// require root rights for the installation run
	if os.Geteuid() != 0 {
		return nil, errors.New("need root permissions to install")
	}

	if !user.UserExists(s.serviceName) {
		return nil, errors.New("user does not exist")
	}

	if !utils.FileExists(servicePath(s.serviceName)) {
		contents := unitFile
		if contents == "" {
			contents, err = serviceFileContents(s.serviceName)
			if err != nil {
				return nil, err
			}
		}
		err := os.WriteFile(servicePath(s.serviceName), []byte(contents), 0644)
		if err != nil {
			return nil, err
		}
	}

	s.refreshStatus()
	err = s.reload()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SystemCtl) Stop() (IdleService, error) {
	_, err := utils.Execute("systemctl stop " + s.serviceName)
	if err != nil {
		return nil, err
	}
	s.refreshStatus()
	return s, nil
}
func (s *SystemCtl) Start() (StartedService, error) {
	_, err := utils.Execute("systemctl start " + s.serviceName)
	if err != nil {
		return nil, err
	}
	s.refreshStatus()
	return s, nil
}

func (s *SystemCtl) IsInstalled() bool {
	return utils.FileExists("/etc/systemd/system/" + s.serviceName + ".service")
}

func (s *SystemCtl) IsRunning() bool {
	return s.status.ActiveState == "active" && s.status.SubState == "running"
}

func Get(name string) (*SystemCtl, error) {
	// require systemctl to be installed
	cmd := exec.Command("which", "systemctl")
	err := cmd.Run()
	if err != nil {
		return nil, errors.New("systemctl is not installed")
	}

	systemCtl := &SystemCtl{
		serviceName: name,
		status:      &ServiceStatus{},
	}

	systemCtl.refreshStatus()

	return systemCtl, nil
}

package systemctl_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/digiconvent/install_on_debian/binary"
	"github.com/digiconvent/install_on_debian/systemctl"
	"github.com/digiconvent/install_on_debian/user"
)

func Main() {
	fmt.Println("MAIN PROGRAM START")
	fmt.Println(os.Args)
	time.Sleep(10 * time.Second)
	fmt.Println("MAIN PROGRAM END")
	os.Exit(0)
}

const name string = "systemctl_tests"

func TestSystemCtl(t *testing.T) {
	defer Cleanup()
	user, err := user.CreateOrGetUser(name)
	if err != nil {
		t.Fatal("Expected to create user for this step", err)
	}
	defer user.Delete()
	sysCtl, err := systemctl.Get(name)
	if err != nil {
		t.Fatal(err)
	}

	// copy the binary to the homefolder of name
	err = binary.New(name).HardLinkToHome()
	if err != nil {
		t.Fatal(err)
	}

	installedService, err := sysCtl.Install("") // use default service file
	if err != nil {
		t.Fatal(err)
	}
	if sysCtl.IsInstalled() == false {
		t.Fatal("expected service to be registered")
		t.FailNow()
	}

	_, err = installedService.Start()
	if err != nil {
		t.Fatal(err)
	}

	_, err = installedService.Uninstall()
	if err != nil {
		t.Fatal(err)
	}
	if sysCtl.IsInstalled() == true {
		t.Fatal("expected IsRegistered() to return false")
	}
}

func Cleanup() {
	sysCtl, _ := systemctl.Get(name)
	if sysCtl == nil {
		return
	}

	sysCtl.Uninstall()
	if sysCtl.User != nil {
		err := sysCtl.User.Delete()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		u, _ := user.CreateOrGetUser(name)
		if u != nil {
			u.Delete()
		}
	}
}

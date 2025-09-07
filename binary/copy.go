package binary

import (
	"os"
	"path"

	"github.com/digiconvent/install_on_debian/utils"
)

// this packages makes sure that the binary is in place and handles "it"

type BinaryOperations interface {
	uri() string
	HardLinkToHome() error
}

type Binary struct {
	name string
}

func (b *Binary) uri() string {
	u, err := os.Executable()
	if err != nil {
		return ""
	}
	return u
}

func (b *Binary) HardLinkToHome() error {
	target := path.Join("/home", b.name, "main")
	if utils.FileExists(target) {
		err := os.Remove(target)
		if err != nil {
			return err
		}
		return b.HardLinkToHome()
	}

	return os.Link(b.uri(), target)
}

func New(name string) BinaryOperations {
	return &Binary{name: name}
}

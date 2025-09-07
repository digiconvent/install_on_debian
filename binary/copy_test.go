package binary_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"
	"testing"

	"github.com/digiconvent/install_on_debian/binary"
	"github.com/digiconvent/install_on_debian/user"
	"github.com/digiconvent/install_on_debian/utils"
)

var name string = "binary_tests"

func TestHardLink(t *testing.T) {
	u, err := user.CreateOrGetUser(name)
	if err != nil {
		t.Fatal(err)
	}

	defer u.Delete()

	if !utils.FileExists("/home/" + name) {
		t.Fatal("for some reason, home folder for", name, "does not exist")
	}

	targetBinaryPath := path.Join("/home", name, "main")
	if utils.FileExists(targetBinaryPath) {
		os.Remove(targetBinaryPath)
	}

	DescribeBinaryStatus(name)
	err = binary.New(name).HardLinkToHome()
	if err != nil {
		t.Fatal(err)
	}
	DescribeBinaryStatus(name)
}

func DescribeBinaryStatus(name string) bool {
	filename := path.Join("/home", name, "main")

	if !utils.FileExists(filename) {
		fmt.Println(filename, "does not exist and it should")
		return false
	}

	fi, err := os.Lstat(filename)
	if err != nil {
		fmt.Println(err)
	}

	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		err = errors.New("cannot convert stat value to syscall.Stat_t")
		fmt.Println(err)
	}

	inode := uint64(s.Ino)
	nlink := uint32(s.Nlink)

	if fi.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(fi.Name())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v is a symlink to %v on inode %v.\n", filename, link, inode)
		os.Exit(0)
	}

	fmt.Printf("The inode for %v, %v, has %v hardlinks.\n", filename, inode, nlink)
	if nlink > 1 {
		fmt.Printf("Inode %v has %v other hardlinks besides %v.\n", inode, nlink, filename)
	} else {
		fmt.Printf("%v is the only hardlink to inode %v.\n", filename, inode)
	}

	return true
}

package user_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/digiconvent/install_on_debian/user"
)

const username string = "someuser"

func TestCreateUser(t *testing.T) {
	defer Cleanup()

	if user.UserExists(username) {
		t.Fatal(username, "should not exist")
	}
	if user.HomeFolderExists(username) {
		t.Fatal("home folder of", username, "should not exist")
	}

	u, err := user.CreateOrGetUser(username)
	if err != nil {
		t.Fatal(err)
	}
	if u == nil {
		t.Fatal("expected a new user")
	}
	// sudoers file exists
	_, err = os.Stat("/etc/sudoers.d/" + username)
	if err != nil {
		t.Fatal(err)
	}

	err = u.Delete()
	if err != nil {
		t.Fatal("could not delete user: " + err.Error())
	}

	// sudoers file should not exist anymore
	_, err = os.Stat("/etc/sudoers.d/" + username)
	if err == nil {
		t.Fatal("expected err not to be nil since file should not exist anymore")
	}
}

func Cleanup() {
	u, err := user.CreateOrGetUser(username)
	if err != nil {
		fmt.Println(err)
	}

	err = u.Delete()
	if err != nil {
		fmt.Println(err)
	}

	if user.UserExists(username) {
		fmt.Println("USER SHOULD NOT EXIST AFTER CLEANUP")
	}
	if user.HomeFolderExists(username) {
		fmt.Println("USER HOME FOLDER SHOULD NOT EXIST AFTER CLEANUP")
	}
}

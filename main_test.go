package install_on_debian_test

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/digiconvent/install_on_debian"
	"github.com/digiconvent/install_on_debian/utils"
)

func main() {
	fmt.Println("Start of the program")
	for i := range 10 {
		time.Sleep(time.Second)
		fmt.Println(10 - i)
	}
	fmt.Println("End of the program")
	os.Exit(0)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestInstallOnDebian(t *testing.T) {
	var name string = randSeq(10)

	fmt.Println("Doing test with " + name)
	if len(os.Args) > 2 {
		fmt.Println("Genuine test")
	} else {
		fmt.Println("Running once and then self destruct")
		main()
	}
	bin := install_on_debian.NewBinary(name)
	service, err := bin.Install()
	if err != nil {
		t.Fatal(err)
	}

	Logs(name, 14)
	time.Sleep(15 * time.Second)
	Logs(name, 14)

	idleService, err := service.Stop()
	if err != nil {
		t.Fatal(err)
	}
	uninstalledService, err := idleService.Uninstall()
	if err != nil {
		t.Fatal(err)
	}

	uninstalledService.DeleteAccount()
}

func Logs(name string, numLines int) {
	c := "journalctl -u " + name + " --reverse"
	output, err := utils.Execute(c)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(output, "\n")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(time.Now().Format(time.RFC1123), c)
	showLines := []string{}
	for i := range numLines {
		if i >= len(lines) {
			break
		}

		lineSegments := strings.Split(lines[i], ": ")
		line := lineSegments[len(lineSegments)-1]
		showLines = append([]string{line}, showLines...)
	}

	fmt.Println(strings.Join(showLines, "\n"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(strings.Repeat("-", 80))
}

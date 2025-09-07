package utils

import (
	"errors"
	"os/exec"
	"strings"
)

func Execute(c string) (string, error) {
	segments := strings.Split(c, " ")
	s, e := exec.Command("sudo", segments...).CombinedOutput()
	if e != nil {
		return "", errors.New(e.Error() + ": " + string(s))
	}
	return string(s), e
}

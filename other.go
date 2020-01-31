package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// UnameMachine return the `uname -m"
func UnameMachine() string {
	var uname syscall.Utsname
	syscall.Uname(&uname)

	arr := uname.Machine[:]
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

// RunCommand run `prog` with `args` and return stdout and stderr combine
func RunCommand(prog string, args ...string) (string, error) {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command(prog, args...)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// RunUrl fetch `url` content.
// NOTE, `args` is not working yet.
func RunUrl(url string, args ...string) (string, error) {
	in, err := OpenURL(url, args...)
	if err != nil {
		return "", err
	}
	defer in.Close()

	bs, err := ioutil.ReadAll(in)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

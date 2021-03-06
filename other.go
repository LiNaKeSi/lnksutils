package lnksutils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
)

// RunCommand run `prog` with `args` and return stdout and stderr combine
func RunCommand(prog string, args ...string) (string, error) {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command(prog, args...)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%q with %v", buf.String(), err)
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

// GetFreePort find a free TCP listen port on the ip
func GetFreePort(ip string) (int, error) {
	if ip == "" {
		ip = "127.0.0.1"
	}
	addr, err := net.ResolveTCPAddr("tcp", ip+":0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func HashFile(filePath string, alg hash.Hash) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(alg, f)
	return hex.EncodeToString(alg.Sum(nil)), err
}

func HASHSelf(alg hash.Hash) (string, error) {
	prog, err := os.Executable()
	if err != nil {
		return "", err
	}
	return HashFile(prog, alg)
}

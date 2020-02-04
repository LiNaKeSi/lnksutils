package utils

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func IsRootPower() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return "0" == currentUser.Uid
}

// LiftPrivilege sudo itself.
//
// Note, this function need zenity and sudo program in system.
func LiftPrivilege() error {
	if IsRootPower() {
		return nil
	}
	password, err := RunCommand("zenity",
		"--password",
		"--title", "权限",
		"--text", "需要特权执行安装过程")
	if err != nil {
		return err
	}
	selfPath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}
	// TODO 避免 PS 看到密码
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("echo %s | sudo -S %s %s",
			password, selfPath, strings.Join(os.Args[1:], " ")))
	cmd.Run()
	os.Exit(0)
	return nil
}

package lnksutils

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
		fmt.Println("无法通过图形界面获取密码，尝试使用sudo执行", err)
	}

	selfPath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if password != "" {
		// TODO 避免 PS 看到密码
		cmd = exec.Command("sh", "-c",
			fmt.Sprintf("echo %s | sudo -E -S %s %s",
				password, selfPath, strings.Join(os.Args[1:], " ")))
	} else {
		cmd = exec.Command("sudo", "-E", selfPath)
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	os.Exit(0)
	return nil
}

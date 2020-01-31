package utils

import (
	"os"
	"os/exec"
	"os/user"
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
func LiftPrivilege(exit bool) error {
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
	err = exec.Command("sh", "-c", "echo "+password+" | sudo -S "+selfPath+"").Start()
	if err != nil {
		return err
	}
	if exit {
		os.Exit(0)
	}
	return nil
}

package liftp

import (
	"bytes"
	"errors"
	"runtime"

	"github.com/gen2brain/dlgs"

	"fmt"
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
func LiftPrivilege(why string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
	default:
		return errors.New("not support")
	}

	if IsRootPower() {
		return nil
	}

	password, ok, err := dlgs.Password("需要特权执行相关操作", why)
	if err != nil {
		fmt.Println("无法通过图形界面获取密码，尝试使用sudo执行", err)
	}
	if !ok {
		return errors.New("用户取消密码输入")
	}

	selfPath, err := os.Executable()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	if password != "" {
		args := []string{
			"-E", "-S", selfPath,
		}
		args = append(args, os.Args[1:]...)
		cmd = exec.Command("sudo", args...)
		cmd.Stdin = bytes.NewReader([]byte(password + string([]byte{'\n', 0})))
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

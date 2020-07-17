//+build linux

package autostart

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/linakesi/lnksutils"
)

type SystemdConf struct {
	Name        string
	Description string
	Exec        string
	WantedBy    string
	IsSession   bool
}

func (conf SystemdConf) InstallFrom(exec string) error {
	err := conf.installExec(exec)
	if err != nil {
		return err
	}

	content, err := conf.buildServiceConf()
	if err != nil {
		return err
	}

	err = lnksutils.SaveToFile(bytes.NewBuffer(content), conf.serviceFile())
	if err != nil {
		return err
	}
	return conf.enableService()
}

func (conf SystemdConf) buildServiceConf() ([]byte, error) {
	const tplSession = `[Unit]
Description={{ .Description }}
After=network.target
[Service]
ExecStart={{ .Exec }}

{{ with .WantedBy }}
[Install]
WantedBy={{ . }}
{{ end }}
`
	const tplSystem = `[Unit]
Description={{ .Description }}
After=network.target

[Service]
Type=simple
ExecStart={{ .Exec }}
Restart=always
RestartSec=5
StartLimitInterval=3
RestartPreventExitStatus=SIGKILL

{{ with .WantedBy }}
[Install]
WantedBy={{ . }}
{{ end }}
`
	tpl := tplSystem
	if conf.IsSession {
		tpl = tplSession
	}
	var buf = bytes.NewBuffer(nil)
	err := template.Must(template.New("eval.systemd.conf").Parse(tpl)).
		Execute(buf, conf)
	return buf.Bytes(), err
}

func (conf SystemdConf) installExec(from string) error {
	if from == "" {
		return nil
	}
	to := strings.Split(conf.Exec, " ")[0]
	if to == "" {
		return fmt.Errorf("There must have valid Exec field in %+v", conf)
	}

	//We must stop the service before installing the program,
	//otherwise it may not be able to successfully overwrite the program.
	conf.ctrl("stop", conf.Name)
	err := lnksutils.FetchFileTo(from, to)
	if err != nil {
		return err
	}
	return os.Chmod(to, 0755)
}

func (conf SystemdConf) ctrl(args ...string) error {
	cmd := exec.Command("systemctl")
	if conf.IsSession {
		cmd.Args = append(cmd.Args, "--user")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append(cmd.Args, args...)
	return cmd.Run()
}

func (conf SystemdConf) serviceFile() string {
	output := "/usr/lib/systemd/system/"
	if conf.IsSession {
		output = os.ExpandEnv("$HOME/.config/systemd/user/")
	}
	return filepath.Join(output, conf.Name+".service")
}

func (conf SystemdConf) enableService() error {
	conf.ctrl("daemon-reload")
	conf.ctrl("enable", conf.Name)
	conf.ctrl("restart", conf.Name)
	return nil
}

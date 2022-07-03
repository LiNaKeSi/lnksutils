package lnksutils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"

	"sync"

	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
)

var initSink sync.Once
var lnkslogSchema = "lnks-log"

// Set the remote log server address,
// the simplest way is run `nc -k -l 7777` on server.
// and invoke SetLogServer("http://$server_ip:7777")
func SetLogServer(addr string) error {
	initSink.Do(func() {
		zap.RegisterSink(lnkslogSchema, newLNKSSink)
	})

	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	u.Scheme = lnkslogSchema

	cfg := logging.GetConfig()
	cfg.URL = u.String()
	logging.SetupLogging(cfg)
	return nil
}

type lnksLogSink struct {
	server string
}

func newLNKSSink(u *url.URL) (zap.Sink, error) {
	return &lnksLogSink{server: "http://" + u.Host + u.Path}, nil
}

func (m *lnksLogSink) Close() error { return nil }
func (m *lnksLogSink) Sync() error  { return nil }

func (m *lnksLogSink) Write(p []byte) (n int, err error) {
	resp, err := http.Post(m.server, "application/binary", bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return len(p), nil
}

func Logger(system string) *logging.ZapEventLogger {
	return logging.Logger(system)
}

func EnableLogDetail(m string) {
	if m == "?" {
		ms := logging.GetSubsystems()
		sort.Strings(ms)
		fmt.Println("All Modules:", ms)
		os.Exit(0)
	} else if m == "warn" {
		logging.SetAllLoggers(logging.LevelWarn)
	} else {
		err := logging.SetLogLevelRegex(m, "debug")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}

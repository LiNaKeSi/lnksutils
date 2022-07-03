package lnksutils

import (
	"fmt"
	"os"
	"sort"

	logging "github.com/ipfs/go-log/v2"
)

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

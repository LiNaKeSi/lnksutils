package lnksutils

import (
	"fmt"
	"runtime"
)

func TraceError(err error, skip int) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return err
	}
	return fmt.Errorf("%s:%d: %s", file, line, err)
}

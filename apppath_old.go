//+build !go1.13

package lnksutils

import "os"

func userCacheDir() (string, error) {
	return os.ExpandEnv("$HOME/.cache"), nil
}

func userConfigDir() (string, error) {
	return os.ExpandEnv("$HOME/.config"), nil
}

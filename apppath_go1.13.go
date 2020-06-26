//+build go1.13

package lnksutils

import "os"

func userCacheDir() (string, error) {
	return os.UserCacheDir()
}

func userConfigDir() (string, error) {
	return os.UserConfigDir()
}

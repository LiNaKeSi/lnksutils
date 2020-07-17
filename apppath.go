package lnksutils

import "path/filepath"

// // AppDir support join user cache/config directory with
// // appname and fields.
type AppDir struct {
	cacheDir  string
	configDir string
}

func (a AppDir) ConfigPath(fields ...string) string {
	return filepath.Join(append([]string{a.configDir}, fields...)...)
}
func (a AppDir) CachePath(fields ...string) string {
	return filepath.Join(append([]string{a.cacheDir}, fields...)...)
}

func NewAppDir(appname string) (AppDir, error) {
	d, err := userCacheDir()
	if err != nil {
		return AppDir{}, err
	}
	c, err := userConfigDir()
	if err != nil {
		return AppDir{}, err
	}
	return AppDir{
		cacheDir:  filepath.Join(d, appname),
		configDir: filepath.Join(c, appname),
	}, nil
}

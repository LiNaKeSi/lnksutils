package lnksutils

import "path"

// // AppDir support join user cache/config directory with
// // appname and fields.
type AppDir struct {
	cacheDir  string
	configDir string
}

func (a AppDir) ConfigPath(fields ...string) string {
	return path.Join(append([]string{a.configDir}, fields...)...)
}
func (a AppDir) CachePath(fields ...string) string {
	return path.Join(append([]string{a.cacheDir}, fields...)...)
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
		cacheDir:  path.Join(d, appname),
		configDir: path.Join(c, appname),
	}, nil
}

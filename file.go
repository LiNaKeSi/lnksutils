package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// OpenURL open the url for reading
// It will reaturn error if open failed or the
// StatusCode is bigger than 299
// NOTE: the return reader need be closed
func OpenURL(url string, args ...string) (io.ReadCloser, error) {
	if len(args) > 0 {
		panic("Doesn't support multiple arguments")
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("OpenURL %q failed %q", url, resp.Status)
	}
	return resp.Body, nil
}

// FetchFile from url.
// Support http protocol and local file
func FetchFile(url string, handle func(r io.Reader) error) error {
	switch {
	case strings.HasPrefix(url, "http"):
		rc, err := OpenURL(url)
		if err != nil {
			return err
		}
		defer rc.Close()
		return handle(rc)
	default:
		f, err := os.Open(url)
		if err != nil {
			return err
		}
		defer f.Close()
		handle(f)
	}
	panic("not reached")
}

// EnsureBaseDir make sure the parent directory of fpath exists
func EnsureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}

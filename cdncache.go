package cdncache

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type cdn struct {
	mountPoint string
	cached     map[string]struct{}
}

var CACHE_DIR = os.Getenv("CDN_CACHE_DIR")

func setCacheDir() {
	if CACHE_DIR == "" {
		CACHE_DIR = filepath.Join("tmp", "cdn_cache")
	}
}

func init() {
	setCacheDir()
}

// CDN returns a function that may be used to return either the cdnUrl
// or a local relative url, depending on mountPoint.
//
// If mountPoint is empty, nothing is cached and the cdnUrl is returned.
//
// Otherwise a fileserver is mounted at mountPoint, serving the files
// from cache directory located at CDN_CACHE_DIR or - if not set - at tmp/cdn_cache.
// Any call to the returned function will then ensure that the corresponding file
// was downloaded to the cache directory and thus can be served by the fileserver.
//
// mountPoint must either be empty or begin and end with /
//
// An invalid mountPoint or an invalid cdnUrl results in a panic.
func CDN(mountPoint string) func(cdnUrl string) string {
	if mountPoint != "" {
		if mountPoint[0] != '/' || mountPoint[len(mountPoint)-1] != '/' {
			panic(`mountpoint must start and end with "/"`)
		}
		http.Handle(mountPoint, http.StripPrefix(mountPoint, http.FileServer(http.Dir(CACHE_DIR))))
	}
	c := &cdn{
		cached:     map[string]struct{}{},
		mountPoint: mountPoint,
	}
	return c.path
}

type Muxer interface {
	Handle(mountpoint string, h http.Handler)
}

func MuxCDN(m Muxer, mountPoint string) func(cdnUrl string) string {
	if mountPoint != "" {
		if mountPoint[0] != '/' || mountPoint[len(mountPoint)-1] != '/' {
			panic(`mountpoint must start and end with "/"`)
		}
		m.Handle(mountPoint, http.StripPrefix(mountPoint, http.FileServer(http.Dir(CACHE_DIR))))
	}
	c := &cdn{
		cached:     map[string]struct{}{},
		mountPoint: mountPoint,
	}
	return c.path
}

// getFile gets the file at cdnUrl and saves it in the CACHE_DIR
func (c *cdn) getFile(cdnUrl string, parsedUrl *url.URL) error {
	file := filepath.Join(CACHE_DIR, parsedUrl.Host, parsedUrl.Path)
	info, err := os.Stat(file)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("%s is a directory", file)
		}
		return nil
	}

	os.MkdirAll(path.Dir(file), 0755)

	if cdnUrl[0:2] == "//" {
		cdnUrl = "http:" + cdnUrl
	}

	resp, errGet := http.Get(cdnUrl)
	if errGet != nil {
		return errGet
	}

	f, errCreate := os.Create(file)
	if errCreate != nil {
		return errCreate
	}

	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// path returns either cdnUrl (no caching) or a corresponding local url
func (c *cdn) path(cdnUrl string) string {
	if c.mountPoint == "" {
		return cdnUrl
	}
	parsedUrl, err := url.Parse(cdnUrl)

	if err != nil {
		panic(fmt.Sprintf("invalid cdnUrl: %s: %s", cdnUrl, err.Error()))
	}

	complete := parsedUrl.Host + parsedUrl.Path

	if _, has := c.cached[complete]; !has {
		err = c.getFile(cdnUrl, parsedUrl)
		if err != nil {
			panic(fmt.Sprintf("can't get cdnUrl: %s: %s", cdnUrl, err.Error()))
		}
		c.cached[complete] = struct{}{}
	}
	return path.Join(c.mountPoint, parsedUrl.Host, parsedUrl.Path)
}

package common

import (
	"io/ioutil"
	"os"
	"path"
)

const (
	CacheDirPermission  = 0777
	CacheFilePermission = 0666
)

func LoadCachedFile(fileName string) ([]byte, error) {
	initFileCache()

	return ioutil.ReadFile(path.Join(_appSettings.CacheDir, fileName))
}

func CacheToFile(fileName string, bytes []byte) error {
	initFileCache()

	return ioutil.WriteFile(path.Join(_appSettings.CacheDir, fileName), bytes, CacheFilePermission)
}

func initFileCache() {
	if _, err := os.Stat(_appSettings.CacheDir); err == nil {
		return
	}

	err := os.MkdirAll(_appSettings.CacheDir, CacheDirPermission)
	if err != nil {
		panic(err)
	}
}

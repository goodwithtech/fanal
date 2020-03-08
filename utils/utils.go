package utils

import (
	"archive/tar"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goodwithtech/deckoder/types"
)

var (
	PathSeparator = fmt.Sprintf("%c", os.PathSeparator)
)

func CacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	return cacheDir
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsCommandAvailable(name string) bool {
	if _, err := exec.LookPath(name); err != nil {
		return false
	}
	return true
}

func IsGzip(f *bufio.Reader) bool {
	buf := make([]byte, 3)
	buf, err := f.Peek(3)
	if err != nil {
		return false
	}
	return buf[0] == 0x1F && buf[1] == 0x8B && buf[2] == 0x8
}

// CreateFilterPathFunc :
func CreateFilterPathFunc(filenames []string) types.FilterFunc {
	return func(h *tar.Header) (bool, error) {
		filePath := filepath.Clean(h.Name)
		fileName := filepath.Base(filePath)
		return StringInSlice(filePath, filenames) || StringInSlice(fileName, filenames), nil
	}
}

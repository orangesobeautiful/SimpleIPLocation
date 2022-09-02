package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// GetEXEDir 取得當前執行檔的目錄路徑
func GetEXEDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	edir := filepath.Dir(exePath)
	return edir, nil
}

func CheckPathExist(fPath string) (bool, error) {
	_, err := os.Stat(fPath)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

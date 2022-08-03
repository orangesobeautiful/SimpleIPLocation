package utils

import (
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

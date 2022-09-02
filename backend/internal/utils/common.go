package utils

import (
	"errors"
	"os"
	"path/filepath"
)

var exeDir string

// InitEXEDirValue 初始化執行檔目錄
func InitEXEDirValue() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir = filepath.Dir(exePath)
	return nil
}

// GetEXEDir 取得當前執行檔的目錄路徑
func GetEXEDir() string {
	return exeDir
}

// PathExist 檢查路徑是否存在
func PathExist(fPath string) (bool, error) {
	_, err := os.Stat(fPath)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

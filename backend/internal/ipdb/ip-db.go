package ipdb

import (
	"errors"
	"path/filepath"

	"SimpleIPLocation/internal/utils"

	"github.com/oschwald/geoip2-golang"
)

var mmdb *geoip2.Reader

// InitIPDB 初始化 IP DB，使用 internal/ipdb 前需要先進行初始化
func InitIPDB() error {
	var err error
	exeDir, err := utils.GetEXEDir()
	if err != nil {
		return errors.New("getEXEDir failed, err=" + err.Error())
	}

	mmdb, err = geoip2.Open(filepath.Join(exeDir, "server-data", "ipdb", "ipdb.mmdb"))
	if err != nil {
		return errors.New("geoip2.Open failed, err=" + err.Error())
	}

	return nil
}

// CloseDB() close ip db, returns the resources to the system.
func CloseDB() error {
	if mmdb == nil {
		return errors.New("mmdb is nil")
	}

	return mmdb.Close()
}

// GetReader 取得 ip db 的 reader
func GetReader() *geoip2.Reader {
	return mmdb
}

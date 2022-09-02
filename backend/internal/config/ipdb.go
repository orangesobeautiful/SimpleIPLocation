package config

import (
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
)

// IPDBConfigInfo IPDB 的設定資料
type IPDBConfigInfo struct {
	Type           string  // 使用的 ip db 類型
	AutoUpdate     bool    // 是否要自動更新資料庫\
	DownSpeedLimit float64 // 下載時的速率限制(bytes/s)

	DBIP DBIPInfo // DB-IP 的相關設定資訊
}

type DBIPInfo struct {
	UpdateDay int // 每月幾號自動更新
}

func parseIPDBConfigFile(ipdbConfig *IPDBConfigInfo) error {
	var err error
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Join("server-data", "config"))
	viper.SetConfigName("ipdb.toml")

	// 讀取設定
	if err = viper.ReadInConfig(); err != nil {
		return errors.New("viper.ReadInConfig err=" + err.Error())
	}

	// 解析設定到 serverConfig
	err = viper.Unmarshal(ipdbConfig)
	return err
}

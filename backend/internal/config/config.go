package config

import (
	"errors"
	"flag"

	"github.com/spf13/viper"
)

// ConfigInfo 伺服器的設定資料
type ConfigInfo struct {
	Host        string // 伺服器要 listen 的 host
	Port        int    // 伺服器要 listen 的 port
	LogFilePath string // log 檔案的路徑
	STDOUT      bool   // 是否要 stdout 輸出 log
	Debug       bool   // 是否為 debug 模式
}

var configInfo ConfigInfo

// InitConfig 初始化設定
func InitConfig() error {
	// 解析設定檔
	if err := parseConfigFile(); err != nil {
		return err
	}

	// 解析執行命令參數
	parseCMD()

	return nil
}

func parseConfigFile() error {
	var err error
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.toml")

	// 讀取設定
	if err = viper.ReadInConfig(); err != nil {
		return errors.New("viper.ReadInConfig err=" + err.Error())
	}

	// 解析設定到 configInfo
	err = viper.Unmarshal(&configInfo)
	return err
}

// parseCMD 解析執行程式的參數
func parseCMD() {
	var cliHost, cliLogFilePath string
	var cliPort int
	var cliStdout, cliDebug bool

	flag.StringVar(&cliHost, "host", "", "Listen Host")
	flag.IntVar(&cliPort, "port", 0, "Listen Port")
	flag.StringVar(&cliLogFilePath, "log", "", "Log file path")
	flag.BoolVar(&cliStdout, "stdout", false, "Stdout Log?")
	flag.BoolVar(&cliDebug, "debug", false, "Debug Mode")

	flag.Parse()

	if cliHost != "" {
		configInfo.Host = cliHost
	}
	if cliPort != 0 {
		configInfo.Port = cliPort
	}
	if cliLogFilePath != "" {
		configInfo.LogFilePath = cliLogFilePath
	}
	if cliStdout {
		configInfo.STDOUT = cliStdout
	}
	if cliDebug {
		configInfo.Debug = cliDebug
	}
}

// GetConfigInfo 取得 config info
func GetConfigInfo() ConfigInfo {
	return configInfo
}

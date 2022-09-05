package config

import (
	"SimpleIPLocation/internal/utils"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
)

var serverConfig ServerConfigInfo
var ipdbConfig IPDBConfigInfo

// InitConfig 初始化設定
func InitConfig() error {
	// 解析設定檔
	configDirPath := filepath.Join(utils.GetEXEDir(), "server-data", "config")
	if dExist, _ := utils.PathExist(configDirPath); !dExist {
		err := os.MkdirAll(configDirPath, utils.NormalDirPerm)
		if err != nil {
			return errors.New("create config dir failed, err=" + err.Error())
		}
	}

	if err := parseServerConfigFile(configDirPath, &serverConfig); err != nil {
		log.Printf("read server config failed, err=%s\n", err)
		log.Printf("using default server config value\n")
		serverConfig.SetToDefault()
	}
	if err := parseIPDBConfigFile(configDirPath, &ipdbConfig); err != nil {
		log.Printf("read ip db config failed, err=%s\n", err)
		log.Printf("using default ip db config value\n")
		ipdbConfig.SetToDefault()
	}

	// 解析執行命令參數
	parseCMD()

	return nil
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
		serverConfig.Host = cliHost
	}
	if cliPort != 0 {
		serverConfig.Port = cliPort
	}
	if cliLogFilePath != "" {
		serverConfig.LogFilePath = cliLogFilePath
	}
	if cliStdout {
		serverConfig.STDOUT = cliStdout
	}
	if cliDebug {
		serverConfig.Debug = cliDebug
	}
}

// GetServerConfig 取得 server config
func GetServerConfig() ServerConfigInfo {
	return serverConfig
}

// GetIPDBConfig 取得 ipdb config
func GetIPDBConfig() IPDBConfigInfo {
	return ipdbConfig
}

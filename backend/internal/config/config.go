package config

import (
	"flag"
)

var serverConfig ServerConfigInfo
var ipdbConfig IPDBConfigInfo

// InitConfig 初始化設定
func InitConfig() error {
	// 解析設定檔
	if err := parseServerConfigFile(&serverConfig); err != nil {
		return err
	}
	if err := parseIPDBConfigFile(&ipdbConfig); err != nil {
		return err
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

package config

import (
	"SimpleIPLocation/internal/utils"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

// ConfigType config info interface
type ConfigType interface {
	SetToDefault()
}

var serverConfig ServerConfigInfo
var ipdbConfig IPDBConfigInfo

func parseConfigFile(configName, configDirPath string, inputConfig ConfigType) error {
	var err error
	v := viper.New()
	v.SetConfigType("toml")
	v.AddConfigPath(configDirPath)
	v.SetConfigName(configName)

	// 讀取設定
	if err = v.ReadInConfig(); err != nil {
		log.Printf("Ignore: read %s config failed, err=%s\n", configName, err)

		// 寫入預設設定
		inputConfig.SetToDefault()
		_ = v.MergeConfigMap(structs.Map(inputConfig))
		writeErr := v.SafeWriteConfig()
		if writeErr != nil {
			err = fmt.Errorf("write default %s config file failed, err=%s", configName, writeErr.Error())
			return err
		}

		log.Printf("using default %s config value\n", configName)
		return nil
	}

	// 解析設定到 serverConfig
	err = v.Unmarshal(inputConfig)
	return err
}

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

	type configRead struct {
		ConfigName    string
		ConfigDirPath string
		ConfigStruct  ConfigType
	}

	configReadList := []configRead{
		{"server", configDirPath, &serverConfig},
		{"ipdb", configDirPath, &ipdbConfig},
	}

	for _, cr := range configReadList {
		if err := parseConfigFile(cr.ConfigName, cr.ConfigDirPath, cr.ConfigStruct); err != nil {
			return err
		}
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

package runenv

import (
	"SimpleIPLocation/internal/config"
	"SimpleIPLocation/internal/ipdb"
	"SimpleIPLocation/internal/utils"
	"errors"
	"os"
	"path/filepath"
)

// InitAll 執行所有需要初始化的參數
func InitAll() error {
	var err error
	err = utils.InitEXEDirValue()
	if err != nil {
		return errors.New("init exe dir value failed, err=" + err.Error())
	}
	serverDataPath := filepath.Join(utils.GetEXEDir(), "server-data")
	if dExist, _ := utils.PathExist(serverDataPath); !dExist {
		err = os.MkdirAll(serverDataPath, utils.NormalDirPerm)
		if err != nil {
			return errors.New("create server-data dir failed, err=" + err.Error())
		}
	}
	if err = config.InitConfig(); err != nil {
		return errors.New("init config failed, err=" + err.Error())
	}
	if err = ipdb.InitIPDB(); err != nil {
		return errors.New("init ip db failed, err=" + err.Error())
	}

	return nil
}

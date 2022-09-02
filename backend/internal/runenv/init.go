package runenv

import (
	"SimpleIPLocation/internal/config"
	"SimpleIPLocation/internal/ipdb"
	"errors"
)

// InitAll 執行所有需要初始化的參數
func InitAll() error {
	var err error
	if err = config.InitConfig(); err != nil {
		return errors.New("init config failed, err=" + err.Error())
	}
	if err = ipdb.InitIPDB(); err != nil {
		return errors.New("init ip db failed, err=" + err.Error())
	}

	return nil
}

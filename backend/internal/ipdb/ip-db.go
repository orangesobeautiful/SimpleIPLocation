package ipdb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"SimpleIPLocation/internal/config"
	"SimpleIPLocation/internal/utils"

	"github.com/oschwald/geoip2-golang"
)

var mmdb *geoip2.Reader
var dbMutex sync.RWMutex

type dbUpdateConfigInfo struct {
	dbSRCType DBSRCType

	downSpeedLimit float64 // 下載速率限制(bytes/s)

	oldPath string // 舊 DB 位置
	orgPath string // 原 DB 位置
	newPath string // 新 DB 位置
}

var dbUpdateConf dbUpdateConfigInfo

// DBSRCType DB 來源類型
type DBSRCType string

const (
	DBIP DBSRCType = "dbip"
)

func StrToDBType(s string) (DBSRCType, bool) {
	s = strings.ToLower(s)
	switch s {
	case "dbip", "db-ip":
		return DBIP, true
	default:
		return "", false
	}
}

// InitIPDB 初始化 IP DB，使用 internal/ipdb 前需要先進行初始化
func InitIPDB() error {
	var err error
	exeDir, err := utils.GetEXEDir()
	if err != nil {
		return errors.New("getEXEDir failed, err=" + err.Error())
	}

	ipdbConfig := config.GetIPDBConfig()
	dbType, tVaild := StrToDBType(ipdbConfig.Type)
	if !tVaild {
		return fmt.Errorf("not support ip db type(%s)", ipdbConfig.Type)
	}

	ipdbDataDir := filepath.Join(exeDir, "server-data", "ipdb")
	if dirExist, _ := utils.CheckPathExist(ipdbDataDir); !dirExist {
		err = os.MkdirAll(ipdbDataDir, utils.NormalDirPerm)
		if err != nil {
			return errors.New("create ipdb data dir failed, err=" + err.Error())
		}
	}

	dbUpdateConf.dbSRCType = dbType
	dbUpdateConf.downSpeedLimit = ipdbConfig.DownSpeedLimit
	dbUpdateConf.newPath = filepath.Join(ipdbDataDir, "ipdb-new.mmdb")
	dbUpdateConf.orgPath = filepath.Join(ipdbDataDir, "ipdb.mmdb")
	dbUpdateConf.oldPath = filepath.Join(ipdbDataDir, "ipdb-old.mmdb")

	dbExist, _ := utils.CheckPathExist(dbUpdateConf.orgPath)
	if dbExist {
		mmdb, err = geoip2.Open(dbUpdateConf.orgPath)
		if err != nil {
			return errors.New("geoip2.Open failed, err=" + err.Error())
		}
	} else {
		err = downloadDB(dbType, dbUpdateConf.newPath, 0)
		if err != nil {
			return err
		}

		err = changeIPDB(dbUpdateConf.newPath, dbUpdateConf.orgPath, dbUpdateConf.oldPath)
		if err != nil {
			return err
		}
	}

	if ipdbConfig.AutoUpdate {
		setNextUpdateTimer(0)
	}

	return nil
}

// downloadDB 根據 DB 資料庫來源下載資料庫
func downloadDB(dbType DBSRCType, dstPath string, downSpeedLimit float64) error {
	switch dbType {
	case DBIP:
		return downloadDBIP(dstPath, downSpeedLimit)
	default:
		return fmt.Errorf("not support ip db type(%s)", dbType)
	}
}

// setNextUpdateTimer 設置下次更新的 timer
func setNextUpdateTimer(td time.Duration) {
	var nextUpdateTimeDuration time.Duration
	if td > 0 {
		nextUpdateTimeDuration = td
	} else {
		nextUpdateTimeDuration = nextAutoUpdateDuration(dbUpdateConf.dbSRCType)
	}
	time.AfterFunc(nextUpdateTimeDuration, updateDBHandler)
}

// updateDBHandler 處理更新的程序
func updateDBHandler() {
	var success = false
	defer func() {
		var nextUpdateTimeDuration time.Duration
		if !success {
			nextUpdateTimeDuration = 24 * time.Hour
		}
		setNextUpdateTimer(nextUpdateTimeDuration)
	}()

	err := downloadDB(
		dbUpdateConf.dbSRCType,
		dbUpdateConf.newPath,
		dbUpdateConf.downSpeedLimit)
	if err != nil {
		log.Printf("downloadDB failed, err=%s", err)
		return
	}

	err = changeIPDB(
		dbUpdateConf.newPath,
		dbUpdateConf.orgPath,
		dbUpdateConf.oldPath)
	if err != nil {
		log.Printf("changeIPDB failed, err=%s", err)
		return
	}

	success = true
}

// nextAutoUpdateDuration 取得下次自動更新的時間
func nextAutoUpdateDuration(dbType DBSRCType) time.Duration {
	nowTime := time.Now()
	ipdbConfig := config.GetIPDBConfig()

	switch dbType {
	case DBIP:
		var updateDay = ipdbConfig.DBIP.UpdateDay
		var nextUpdateTime time.Time
		if nowTime.Day() >= ipdbConfig.DBIP.UpdateDay {
			nextUpdateTime = time.Date(nowTime.Year(), nowTime.Month()+1, updateDay, 0, 0, 0, 0, nowTime.Location())
		} else {
			nextUpdateTime = time.Date(nowTime.Year(), nowTime.Month(), updateDay, 0, 0, 0, 0, nowTime.Location())
		}
		return nextUpdateTime.Sub(nowTime)
	default:
		panic("unsupport db type(" + dbType + ")")
	}
}

// changeIPDB 將舊的 ipdb 替換成新的
func changeIPDB(newPath, orgPath, oldPath string) error {
	var err error
	if fExist, _ := utils.CheckPathExist(oldPath); fExist {
		err = os.RemoveAll(oldPath)
		if err != nil {
			return errors.New("remove old ipdb failed, err=" + err.Error())
		}
	}

	if fExist, _ := utils.CheckPathExist(orgPath); fExist {
		err = os.Rename(orgPath, oldPath)
		if err != nil {
			return errors.New("rename ipdb org to old failed, err=" + err.Error())
		}
	}
	err = os.Rename(newPath, orgPath)
	if err != nil {
		return errors.New("rename ipdb new to org failed, err=" + err.Error())
	}

	newMMDB, err := geoip2.Open(orgPath)
	if err != nil {
		return errors.New("geoip2.Open failed, err=" + err.Error())
	}

	dbMutex.Lock()
	_ = CloseDB()
	mmdb = newMMDB
	dbMutex.Unlock()

	return nil
}

// CloseDB close ip db, returns the resources to the system.
func CloseDB() error {
	if mmdb == nil {
		return errors.New("mmdb is nil")
	}

	return mmdb.Close()
}

// GetReader 取得 ip db 的 reader
func GetReader() *geoip2.Reader {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	return mmdb
}

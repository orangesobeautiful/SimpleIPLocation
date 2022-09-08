package config

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

func (c *IPDBConfigInfo) SetToDefault() {
	c.Type = "dbip"
	c.AutoUpdate = true
	c.DownSpeedLimit = 0
	c.DBIP.UpdateDay = 5
}

package config

// ServerConfigInfo 伺服器的設定資料
type ServerConfigInfo struct {
	Host        string // 伺服器要 listen 的 host
	Port        int    // 伺服器要 listen 的 port
	LogFilePath string // log 檔案的路徑
	STDOUT      bool   // 是否要 stdout 輸出 log
	Debug       bool   // 是否為 debug 模式
}

func (c *ServerConfigInfo) SetToDefault() {
	c.Host = "0.0.0.0"
	c.Port = 80
	c.LogFilePath = ""
	c.STDOUT = false
	c.Debug = false
}

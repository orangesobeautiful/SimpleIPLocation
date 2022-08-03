package controllers

import (
	"encoding/json"
	"net/http"
	"net/netip"

	"SimpleIPLocation/internal/ipdb"
	"SimpleIPLocation/internal/utils"

	"github.com/julienschmidt/httprouter"
)

type ipInfoResp struct {
	IP        string  // IP 位置
	Continent string  // 大陸/州
	Country   string  // 國家
	City      string  // 城市
	Longitude float64 // 經度
	Latitude  float64 // 緯度
}

// IPInfo 處理 ip 查詢 IP 資訊的 controller
func IPInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 嘗試提取指定解析的 IP 位置
	var ipStr string
	ipStr = ps.ByName("ip")
	if len(ipStr) > 0 && ipStr[0] == '/' {
		ipStr = ipStr[1:]
	}

	var ip netip.Addr
	if ipStr == "" {
		// 若沒有指定查詢 IP，則根據客戶端請求 IP 查詢
		ip = utils.ParseReqRemoteIP(r)
	} else {
		ip, _ = netip.ParseAddr(ipStr)
	}

	if !ip.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("ip format was error"))
		return
	}

	// 查詢 IP 資料庫中的資訊
	mmdb := ipdb.GetReader()
	record, err := mmdb.City(utils.NetIPAddr2netIP(ip))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("IP Not Found"))
		return
	}

	resp := ipInfoResp{
		IP:        ip.String(),
		Continent: record.Continent.Names["en"],
		Country:   record.Country.IsoCode,
		City:      record.City.Names["en"],
		Longitude: record.Location.Longitude,
		Latitude:  record.Location.Latitude,
	}

	// 回傳 json 格式
	w.Header().Set("content-type", "application/json; charset=utf-8")
	jEncoder := json.NewEncoder(w)
	_ = jEncoder.Encode(resp)
}

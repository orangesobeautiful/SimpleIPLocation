package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"

	"SimpleIPLocation/internal/config"

	"github.com/julienschmidt/httprouter"
	"github.com/oschwald/geoip2-golang"
)

var mmdb *geoip2.Reader

func getEXEDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	edir := filepath.Dir(exePath)
	return edir, nil
}

type ipInfoResp struct {
	IP        string  // IP 位置
	Continent string  // 大陸/州
	Country   string  // 國家
	City      string  // 城市
	Longitude float64 // 經度
	Latitude  float64 // 緯度
}

func netipAddr2netIP(addr netip.Addr) net.IP {
	if addr.Is4() {
		p := make(net.IP, net.IPv6len)
		p[10] = 0xff
		p[11] = 0xff
		b4Ary := addr.As4()
		p[12] = b4Ary[0]
		p[13] = b4Ary[1]
		p[14] = b4Ary[2]
		p[15] = b4Ary[3]
		return p
	}
	return net.IP(addr.AsSlice())
}

func IPInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ipStr string
	var ip netip.Addr
	ipStr = ps.ByName("ip")
	if len(ipStr) > 0 && ipStr[0] == '/' {
		ipStr = ipStr[1:]
	}

	if ipStr == "" {
		ip = ParseReqRemoteIP(r)
	} else {
		ip, _ = netip.ParseAddr(ipStr)
	}

	if !ip.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("ip format was error"))
		return
	}

	record, err := mmdb.City(netipAddr2netIP(ip))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("IP Not Found"))
		return
	}

	resp := ipInfoResp{
		IP:        ipStr,
		Continent: record.Continent.Names["en"],
		Country:   record.Country.IsoCode,
		City:      record.City.Names["en"],
		Longitude: record.Location.Longitude,
		Latitude:  record.Location.Latitude,
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	jEncoder := json.NewEncoder(w)
	_ = jEncoder.Encode(resp)
}

func main() {
	var err error

	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	exeDir, err := getEXEDir()
	if err != nil {
		log.Print("getEXEDir failed, err=", err)
		return
	}

	mmdb, err = geoip2.Open(filepath.Join(exeDir, "server-data", "ipdb", "ipdb.mmdb"))
	if err != nil {
		log.Print(err)
		return
	}

	if err = config.InitConfig(); err != nil {
		log.Print(err)
		return
	}
	configInfo := config.GetConfigInfo()

	router := httprouter.New()
	router.GET("/ipinfo", IPInfo)
	router.GET("/ipinfo/:ip", IPInfo)

	address := fmt.Sprintf("%s:%d", configInfo.Host, configInfo.Port)
	fmt.Printf("Start Listen at %s\n", address)

	if err = http.ListenAndServe(address, router); err != nil {
		log.Print(err)
		return
	}
}

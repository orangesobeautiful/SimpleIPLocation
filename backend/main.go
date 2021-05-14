package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

var mmdb *geoip2.Reader

func PrintGeoip2CityInfo(record *geoip2.City) {
	fmt.Println("大陸 Code:\t", record.Continent.Code)
	fmt.Println("大陸名稱:\t", record.Continent.Names)
	fmt.Println("大陸 GeoNameID:\t", record.Continent.GeoNameID)

	fmt.Println("國家名稱:\t", record.Country.Names)
	fmt.Println("國家 IsoCode:\t", record.Country.IsoCode)
	fmt.Println("國家 GeoNameID:\t", record.Country.GeoNameID)
	fmt.Println("是否為歐盟?:\t", record.Country.IsInEuropeanUnion)

	fmt.Println("城市名稱:\t", record.City.Names)
	fmt.Println("城市 GeoNameID:\t", record.City.GeoNameID)

	fmt.Println("AccuracyRadius:\t", record.Location.AccuracyRadius)
	fmt.Println("經度:\t", record.Location.Longitude)
	fmt.Println("緯度:\t", record.Location.Latitude)
	fmt.Println("MetroCode:\t", record.Location.MetroCode)
	/*
		fmt.Println("時區:\t", record.Location.TimeZone)

		fmt.Println("郵遞區號:\t", record.Postal.Code)

		fmt.Println("RegisteredCountry 名稱:\t", record.RegisteredCountry.Names)
		fmt.Println("RegisteredCountry IsoCode:\t", record.RegisteredCountry.IsoCode)
		fmt.Println("RegisteredCountry GeoNameID:\t", record.RegisteredCountry.GeoNameID)
		fmt.Println("RegisteredCountry 是否為歐盟?:\t", record.RegisteredCountry.IsInEuropeanUnion)

		fmt.Println("record.RepresentedCountry:\t", record.RepresentedCountry)

		fmt.Println("RepresentedCountry 名稱:\t", record.RepresentedCountry.Names)
		fmt.Println("RepresentedCountry IsoCode:\t", record.RepresentedCountry.IsoCode)
		fmt.Println("RepresentedCountry GeoNameID:\t", record.RepresentedCountry.GeoNameID)
		fmt.Println("RepresentedCountry 類型:\t", record.RepresentedCountry.Type)
		fmt.Println("RepresentedCountry 是否為歐盟?:\t", record.RepresentedCountry.IsInEuropeanUnion)

		fmt.Println("Subdivisions:\t", record.Subdivisions)

		fmt.Println("Traits:\t", record.Traits)
	*/
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.Host
}

func IPInfo(c *gin.Context) {
	if c.Request.Method == "GET" {
		var ipStr string
		ipStr = c.Param("ip")
		if len(ipStr) > 0 && ipStr[0] == '/' {
			ipStr = ipStr[1:]
		}
		if ipStr == "" {
			ipStr = GetIP(c.Request)
		}
		ip := net.ParseIP(ipStr)

		if ip == nil {
			c.String(400, "Bad Request")
			return
		}

		record, err := mmdb.City(ip)
		if err != nil {
			c.String(404, "IP Not Found")
			return
		}

		//PrintGeoip2CityInfo(record)

		c.JSON(200, gin.H{"IP": ipStr,
			"Continent":   record.Continent.Names["en"],
			"Country":     record.Country.IsoCode,
			"City":        record.City.Names["en"],
			"Longitude": record.Location.Longitude,
			"Latitude":  record.Location.Latitude,
		})
	}
}

func main() {

	var err error

	mmdb, err = geoip2.Open("dpip/dbip-city-lite-2021-05.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer mmdb.Close()

	var cliHost = flag.String("host", "", "Listen Host")
	var cliPort = flag.String("port", "80", "Listen Port")
	var cliLogFile = flag.String("log", "", "Log file path")
	var cliStdout = flag.Bool("stdout", false, "Stdout Log?")
	var cliDebug = flag.Bool("debug", false, "Debug Mode")

	flag.Parse()

	host := *cliHost
	port := *cliPort
	output := *cliLogFile
	stdout := *cliStdout
	debugMode := *cliDebug

	var f *os.File

	if port == "env" {
		if v := os.Getenv("PORT"); len(v) > 0 {
			port = v
		}
	}

	var ioWriterList []io.Writer

	if stdout {
		gin.ForceConsoleColor()
		ioWriterList = append(ioWriterList, os.Stdout)
	}

	if output != "" {
		f, err = os.OpenFile(output, os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		ioWriterList = append(ioWriterList, f)

		defer f.Close()
	}

	gin.DefaultWriter = io.MultiWriter(ioWriterList...)

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/ipinfo/", IPInfo)
	r.GET("/ipinfo/:ip", IPInfo)

	fmt.Println("Start Listen!")
	err = r.Run(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("End Server")
}

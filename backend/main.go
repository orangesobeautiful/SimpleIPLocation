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
			c.String(http.StatusBadRequest, "Bad Request")
			return
		}

		record, err := mmdb.City(ip)
		if err != nil {
			c.String(http.StatusNotFound, "IP Not Found")
			return
		}

		c.JSON(http.StatusOK, gin.H{"IP": ipStr,
			"Continent": record.Continent.Names["en"],
			"Country":   record.Country.IsoCode,
			"City":      record.City.Names["en"],
			"Longitude": record.Location.Longitude,
			"Latitude":  record.Location.Latitude,
		})
	}
}

func main() {
	var err error

	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	mmdb, err = geoip2.Open("dpip/dbip-city-lite-2021-05.mmdb")
	if err != nil {
		log.Print(err)
		return
	}

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
		const NormalFileMode os.FileMode = 0644
		f, err = os.OpenFile(output, os.O_CREATE|os.O_APPEND, NormalFileMode)
		if err != nil {
			log.Print(err)
			return
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
		log.Print(err)
		return
	}

	fmt.Println("End Server")
}

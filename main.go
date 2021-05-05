package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
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

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {

		ipstr := GetIP(r)

		ip := net.ParseIP(ipstr)

		record, err := mmdb.City(ip)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Your IP: %v<br/>", ipstr)
		fmt.Fprintf(w, "continent: %v<br/>", record.Continent.Names["en"])
		fmt.Fprintf(w, "country: %v<br/>", record.Country.Names["en"])
		fmt.Fprintf(w, "country ISO Code: %v<br/>", record.Country.IsoCode)
		fmt.Fprintf(w, "city: %v<br/>", record.City.Names["en"])
		fmt.Fprintf(w, "Coordinates: %v, %v<br/><br/>", record.Location.Latitude, record.Location.Longitude)
		fmt.Fprint(w, "<a href='https://db-ip.com'>IP Geolocation by DB-IP</a>")

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
	var cliProxyHeader = flag.Bool("proxy", false, "Has Proxy Header?")

	flag.Parse()

	host := *cliHost
	port := *cliPort
	output := *cliLogFile
	isProxyHeader := *cliProxyHeader

	r := httprouter.New()
	r.GET("/", index)

	var f *os.File

	var serverHandler http.Handler

	if port == "env" {
		if v := os.Getenv("PORT"); len(v) > 0 {
			port = v
		}
	}

	if output == "no" {
		serverHandler = r
	} else if output != "" {
		f, err = os.OpenFile(output, os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		serverHandler = handlers.LoggingHandler(f, r)
		serverHandler = handlers.LoggingHandler(os.Stdout, serverHandler)

		defer f.Close()
	} else {
		serverHandler = handlers.LoggingHandler(os.Stdout, r)
	}

	if isProxyHeader {
		serverHandler = handlers.ProxyHeaders(serverHandler)
	}

	fmt.Println("Start Listen!")
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), serverHandler)

	fmt.Println("End Server")
}

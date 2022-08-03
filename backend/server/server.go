package server

import (
	"SimpleIPLocation/controllers"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NewRouter create a new http router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/ipinfo", controllers.IPInfo)
	router.GET("/ipinfo/:ip", controllers.IPInfo)

	return router
}

// RunServer 執行 Server
func RunServer(host string, port int) error {
	router := NewRouter()

	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Start Listen at %s\n", address)

	var err error
	err = http.ListenAndServe(address, router)
	return err
}

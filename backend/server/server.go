package server

import (
	"SimpleIPLocation/controllers"
	"SimpleIPLocation/internal/httpfs"
	"SimpleIPLocation/internal/utils"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// NewRouter create a new http router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/ipinfo", controllers.IPInfo)
	router.GET("/ipinfo/:ip", controllers.IPInfo)

	spaFS := http.Dir(filepath.Join(utils.GetEXEDir(), "server-data", "public", "original"))
	router.NotFound = httpfs.NewFileServer(spaFS, true)
	return router
}

// RunServer 執行 Server
func RunServer(host string, port int) error {
	router := NewRouter()

	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Start Listen at %s\n", address)

	err := http.ListenAndServe(address, router)
	return err
}

package server

import (
	"SimpleIPLocation/controllers"
	"SimpleIPLocation/internal/httpfs"
	"SimpleIPLocation/internal/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

// NewRouter create a new http router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/ipinfo", controllers.IPInfo)
	router.GET("/ipinfo/:ip", controllers.IPInfo)

	spaFS := http.Dir(filepath.Join(utils.GetEXEDir(), "frontend-static"))
	router.NotFound = httpfs.NewFileServer(spaFS, true)
	return router
}

// RunServer 執行 Server
func RunServer(host string, port int) error {
	router := NewRouter()

	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Start Listen at %s\n", address)

	httpServer := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err := httpServer.ListenAndServe()
	return err
}

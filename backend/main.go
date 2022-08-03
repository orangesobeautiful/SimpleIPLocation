package main

import (
	"log"

	"SimpleIPLocation/internal/config"
	"SimpleIPLocation/internal/runenv"
	"SimpleIPLocation/server"
)

func main() {
	var err error

	err = runenv.InitAll()
	if err != nil {
		log.Fatalf("init failed, err=%s", err)
	}

	configInfo := config.GetConfigInfo()

	err = server.RunServer(configInfo.Host, configInfo.Port)
	if err != nil {
		log.Fatalf("run server failed, err=%s", err)
	}
}

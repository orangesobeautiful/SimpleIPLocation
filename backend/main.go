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

	serverConfig := config.GetServerConfig()

	err = server.RunServer(serverConfig.Host, serverConfig.Port)
	if err != nil {
		log.Fatalf("run server failed, err=%s", err)
	}
}

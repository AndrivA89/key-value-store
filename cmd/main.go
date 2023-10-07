package main

import (
	"os"

	"github.com/labstack/gommon/log"

	"github.com/AndrivA89/key-value-store/internal/entity"
	"github.com/AndrivA89/key-value-store/internal/logger"
	"github.com/AndrivA89/key-value-store/internal/server"
	"github.com/AndrivA89/key-value-store/internal/store"
)

func main() {
	var (
		loggerType = entity.LoggerType(os.Getenv("LOGGER_TYPE"))
		serverType = entity.ServerType(os.Getenv("SERVER_TYPE"))
	)

	l, err := logger.NewLogger(loggerType)
	if err != nil {
		log.Fatal(err, loggerType)
	}

	keyValueStore, err := store.NewStore(l)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(serverType)
	err = s.Start(keyValueStore)
	if err != nil {
		log.Fatal(err, serverType)
	}
}

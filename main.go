package main

import (
	"encoding/gob"
	"os"
	"simple-api/handlers"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func main() {
	gob.Register(uuid.UUID{})
	_ = godotenv.Load()
	handler, err := handlers.New()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	handlers.InitRoutes(handler)
	defer handler.DB.Close()
	handler.E.Logger.Info(handler.E.Start(":8080"))
}

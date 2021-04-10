package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/xeynse/XeynseJar_analytics/internal/config"
	"github.com/xeynse/XeynseJar_analytics/internal/handler"
	"github.com/xeynse/XeynseJar_analytics/internal/iot"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/jar"
	"github.com/xeynse/XeynseJar_analytics/internal/resource/db"
	"github.com/xeynse/XeynseJar_analytics/internal/resource/file"
	jarUseCase "github.com/xeynse/XeynseJar_analytics/internal/usecase/jar"
)

func main() {
	fileResource := file.New()
	config, err := config.Init(fileResource)
	if err != nil {
		log.Fatal("[Main] Fatal initializing config :", err, " env :", os.Getenv("XEYNSEENV"))
	}

	dbResource, err := db.New(config)
	if err != nil {
		log.Fatal("[Main] Fatal connecting database :", err)
	}

	jarRepo := jar.New(dbResource)

	iotHub := iot.New(jarRepo)
	go iotHub.Run()

	jarUseCase := jarUseCase.New(config, jarRepo)

	router := httprouter.New()
	handler.New(router, jarUseCase)
	server := &http.Server{
		Handler: router,
		Addr:    config.Server.Address,
	}
	log.Fatal(server.ListenAndServe())
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/xeynse/XeynseJar_analytics/internal/config"
	"github.com/xeynse/XeynseJar_analytics/internal/handler"
	"github.com/xeynse/XeynseJar_analytics/internal/iot"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/home"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/jar"
	dbanalytics "github.com/xeynse/XeynseJar_analytics/internal/resource/db/analytics"
	dbhomeconfig "github.com/xeynse/XeynseJar_analytics/internal/resource/db/homeconfig"
	"github.com/xeynse/XeynseJar_analytics/internal/resource/file"
	jarUseCase "github.com/xeynse/XeynseJar_analytics/internal/usecase/jar"
)

func main() {
	fileResource := file.New()
	config, err := config.Init(fileResource)
	if err != nil {
		log.Fatal("[Main] Fatal initializing config :", err, " env :", os.Getenv("XEYNSEENV"))
	}

	homeconfigDB := dbhomeconfig.New()

	anaylticsDB, err := dbanalytics.New(config)
	if err != nil {
		log.Fatal("[Main] Fatal connecting database :", err)
	}

	jarRepo := jar.New(anaylticsDB)

	iotHub, err := iot.New(jarRepo)
	if err != nil {
		log.Fatal("[Main] Fatal initializing iotHub :", err, " env :", os.Getenv("XEYNSEENV"))
	}
	go iotHub.Run()

	homeconfigRepo := home.New(config, homeconfigDB)
	jarUseCase := jarUseCase.New(config, homeconfigRepo, jarRepo)

	router := httprouter.New()
	handler.New(router, jarUseCase)
	server := &http.Server{
		Handler: router,
		Addr:    config.Server.Address,
	}
	log.Fatal(server.ListenAndServe())
}

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yac-d/sipb/configdef"
	"github.com/yac-d/sipb/filebin/sqlfb"
	"github.com/yac-d/sipb/httpsrv"
	"github.com/yac-d/sipb/logger"
)

func main() {
	var config configdef.Config
	logger.LogConfigRead("./config.yaml", config.ReadFromYAML("./config.yaml"))
	// Overrides config from file only for environment variables that are set (unset ones are ignored)
	logger.LogConfigRead("environment variables", config.ReadFromEnvVars())

	var bin = sqlfb.New(config)
	if err := bin.Initialize(); err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	var srv = httpsrv.New(config, bin)
	srv.OnSave = logger.LogFileSave
	srv.OnDetailsRequested = logger.LogFileDetailsRequest
	srv.OnCountRequested = logger.LogFileCountRequest

	logger.Log(srv.Start())

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator

	bin.Cleanup()
	logger.Log("Exiting")
}
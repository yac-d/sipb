package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin/simplefb"
	"github.com/Eeshaan-rando/sipb/src/httpsrv"
	"github.com/Eeshaan-rando/sipb/src/logger"
)

func main() {
	var config configdef.Config
	logger.LogConfigRead("./config.yaml", config.ReadFromYAML("./config.yaml"))
	// Overrides config from file only for environment variables that are set (unset ones are ignored)
	logger.LogConfigRead("environment variables", config.ReadFromEnvVars())

	var bin = simplefb.New(config)
	var srv = httpsrv.New(config, bin)
	srv.OnSave = logger.LogFileSave
	srv.OnDetailsRequested = logger.LogFileDetailsRequest
	srv.OnCountRequested = logger.LogFileCountRequest

	logger.Log(srv.Start())

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
	logger.Log("Exiting")
}

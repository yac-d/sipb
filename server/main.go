package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin"
	"github.com/Eeshaan-rando/sipb/src/filebin/simplefb"
	"github.com/Eeshaan-rando/sipb/src/httpsrv"
)

func logFileSave(result filebin.SaveFileResult) {
	log.Printf("Saving file %s", result.Filename)
	if result.Error != nil {
		log.Printf("Error while saving file: %s", result.Error)
	}
	if result.TruncatedBytes > 0 {
		log.Printf("File %s above set size limit, truncating by %d bytes", result.Filename, result.TruncatedBytes)
	}
}

func logFileDetailsRequest(result filebin.FileDetailsResult) {
	log.Printf("Details for file %s requested", result.Filename)
	if result.Error != nil {
		log.Printf("Error while getting file details: %s")
	}
}

func logFileCountRequest(result filebin.FileCountResult) {
	if result.Error != nil {
		log.Printf("Error while getting file count: %s", result.Error)
	} else {
		log.Printf("File count requested; currently at %d", result.Count)
	}
}

func main() {
	var config configdef.Config
	if err := config.ReadFromYAML("./config.yaml"); err != nil {
		log.Fatalln("Error reading configuration from ./config.yaml")
	}
	log.Printf("Read configuration from ./config.yaml")

	// Overrides config from file only for environment variables that are set (unset ones are ignored)
	if err := config.ReadFromEnvVars(); err != nil {
		log.Fatalln("Error reading configuration from environment variables")
	}
	log.Printf("Read configuration from environment variables")

	var bin = simplefb.New(config)
	var srv = httpsrv.New(config, bin)
	srv.OnSave = logFileSave
	srv.OnDetailsRequested = logFileDetailsRequest
	srv.OnCountRequested = logFileCountRequest

	log.Println(srv.Start())

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
	log.Println("Exiting")
}

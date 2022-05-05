package main

import (
	"os"
	"os/signal"
	"syscall"
	"log"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin/simplefb"
	"github.com/Eeshaan-rando/sipb/src/httpsrv"
)

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

	log.Println(srv.Start())

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
	log.Println("Exiting")
}

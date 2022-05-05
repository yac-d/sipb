package main

import (
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"log"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	sfb "github.com/Eeshaan-rando/sipb/src/filebin/simplefb"
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

	var bin = sfb.NewFromConfig(config)

	saveFile := func(w http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(64000000)
		incomingFile, h, err := request.FormFile("file")

		log.Printf("Receiving file %s", h.Filename)
		if err != nil {
			log.Println("Error reading uploaded file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if bin.SaveFile(incomingFile, h) { // checking if file had to be truncated
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		}
	}

	retrieveFileDetails := func(w http.ResponseWriter, request *http.Request) {
		incomingLen, _ := strconv.Atoi(request.Header["Content-Length"][0])
		var buf = make([]byte, incomingLen)
		request.Body.Read(buf)
		whichFile, err := strconv.Atoi(string(buf))

		details, err := bin.DetailsOfNthNewest(whichFile)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Invalid request for file details")
			return
		}

		w.Write(details.AsJSON())
		log.Printf("File %s requested", details.Path)
	}

	retrieveFileCount := func(w http.ResponseWriter, request *http.Request) {
		cnt := bin.Count()
		w.Write([]byte(strconv.Itoa(cnt)))
		log.Printf("File count requested, currently at %d", cnt)
	}

	http.Handle("/", http.FileServer(http.Dir(config.WebpageDir)))
	http.HandleFunc("/upload", saveFile)
	http.HandleFunc("/retrieve", retrieveFileDetails)
	http.HandleFunc("/retrieve/fileCount", retrieveFileCount)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.BindAddr, config.Port), nil)
	log.Println(err)

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
	log.Println("Exiting")
}

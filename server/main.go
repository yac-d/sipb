package main

import (
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"io"
	"strconv"
	"log"
	"bytes"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin"
)

func main() {
	var config configdef.Config
	config.ReadFromYAML("./config.yaml")
	log.Printf("Read configuration from ./config.yaml")
	// Overrides config from file only for environment variables that are set (unset ones are ignored)
	config.ReadFromEnvVars()
	log.Printf("Read configuration from environment variables")

	var bin = filebin.NewFromConfig(config)

	saveFile := func(w http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(32000000)
		incomingFile, h, err := request.FormFile("file")

		log.Printf("Receiving file %s", h.Filename)
		if err != nil {
			log.Println("Error reading uploaded file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bin.SaveFile(incomingFile, h)
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

func readerLen(r io.Reader) int {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return len(buf.Bytes())
}

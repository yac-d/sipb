package main

import (
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"io"
	"io/ioutil"
	"strconv"
	"time"
	"log"
	"bytes"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/utils"
	"github.com/Eeshaan-rando/sipb/src/filedetails"
)

func main() {
	var config configdef.Config
	config.ReadFromYAML("./config.yaml")
	log.Printf("Read configuration from ./config.yaml")
	// Overrides config from file only for environment variables that are set (unset ones are ignored)
	config.ReadFromEnvVars()
	log.Printf("Read configuration from environment variables")

	if !utils.FileExists(config.BinDir) {
		os.MkdirAll(config.BinDir, 0755)
		log.Printf("Creating bin directory %s", config.BinDir)
	}

	saveFile := func(w http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(32000000)
		incomingFile, h, err := request.FormFile("file")

		if err != nil {
			log.Println("Error reading uploaded file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var filename = strconv.Itoa(int(time.Now().UnixMilli())) + "_" + h.Filename
		persistedFile, _ := os.Create(path.Join(config.BinDir, filename))
		log.Printf("Receiving file %s", h.Filename)
		if config.MaxFileSize > -1 {
			io.CopyN(persistedFile, incomingFile, config.MaxFileSize)
			l := readerLen(incomingFile)
			if int64(l) > config.MaxFileSize {
				log.Printf("File %s above set size limit, truncating from %d to %d bytes", filename, l, config.MaxFileSize)
			}
		} else {
			io.Copy(persistedFile, incomingFile)
		}
		persistedFile.Close()
		log.Printf("File %s saved as %s", h.Filename, filename)

		if config.MaxFileCnt != -1 {
			files, _ := ioutil.ReadDir(config.BinDir)
			var i = 0
			for len(files) - i > config.MaxFileCnt {
				os.Remove(path.Join(config.BinDir, files[i].Name()))
				i += 1
				log.Printf("Removed old file %s", files[i].Name())
			}
		}
	}

	retrieveFileDetails := func(w http.ResponseWriter, request *http.Request) {
		files, _ := ioutil.ReadDir(config.BinDir)

		incomingLen, _ := strconv.Atoi(request.Header["Content-Length"][0])
		var buf = make([]byte, incomingLen)
		request.Body.Read(buf)
		whichFile, err := strconv.Atoi(string(buf))

		if whichFile > len(files) || whichFile < 1 || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Invalid request for file details")
			return
		}

		var filename = files[len(files) - whichFile].Name()
		var details = filedetails.NewForFile(filename)

		w.Write(details.AsJSON())
		log.Printf("File %s requested", filename)
	}

	retrieveFileCount := func(w http.ResponseWriter, request *http.Request) {
		files, _ := ioutil.ReadDir(config.BinDir)
		w.Write([]byte(strconv.Itoa(len(files))))
		log.Printf("File count requested, currently at %d", len(files))
	}

	filedetails.SetFilesystemLocation(config.BinDir)
	filedetails.SetURLPath(config.BinPath)

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

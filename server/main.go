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
	"encoding/json"
	"log"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/utils"
)

func main() {
	var config configdef.Config
	config.ReadFromYAML("./config.yaml")

	if !utils.FileExists(config.BinDir) {
		os.MkdirAll(config.BinDir, 0755)
	}

	saveFile := func(w http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(32000000)
		incomingFile, h, _ := request.FormFile("file")

		var filename = strconv.Itoa(int(time.Now().UnixMilli())) + "_" + h.Filename
		persistedFile, _ := os.Create(path.Join(config.BinDir, filename))
		io.Copy(persistedFile, incomingFile)
		persistedFile.Close()

		if config.MaxFileCnt != -1 {
			files, _ := ioutil.ReadDir(config.BinDir)
			var i = 0
			for len(files) - i > config.MaxFileCnt {
				os.Remove(path.Join(config.BinDir, files[i].Name()))
				i += 1
			}
		}
	}

	retrieveFileDetails := func(w http.ResponseWriter, request *http.Request) {
		files, _ := ioutil.ReadDir(config.BinDir)

		incomingLen, _ := strconv.Atoi(request.Header["Content-Length"][0])
		var buf = make([]byte, incomingLen)
		request.Body.Read(buf)
		whichFile, _ := strconv.Atoi(string(buf))

		if whichFile > len(files) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var details = make(map[string]string)

		f, _ := os.Open(path.Join(config.BinDir, files[len(files) - whichFile].Name()))
		var fHeader = make([]byte, 512)
		f.Read(fHeader)
		f.Close()

		details["Type"] = http.DetectContentType(fHeader)
		details["Path"] = path.Join(config.BinPath, files[len(files) - whichFile].Name())
		log.Println("Requested:", details)
		outgoing, _ := json.Marshal(details)
		w.Write(outgoing)
	}

	retrieveFileCount := func(w http.ResponseWriter, request *http.Request) {
		files, _ := ioutil.ReadDir(config.BinDir)
		w.Write([]byte(strconv.Itoa(len(files))))
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
}

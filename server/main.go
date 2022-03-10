package main

import (
	"net"
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

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/utils"
)

func getIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return net.IP(conn.LocalAddr().(*net.UDPAddr).IP).String()
}

func main() {
	var config configdef.Config
	config.ReadFromYAML("./config.yaml")

	if !utils.FileExists(config.BinDir) {
		os.MkdirAll(config.BinDir, 0755)
	}

	saveImage := func(w http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(32000000)
		incomingFile, h, _ := request.FormFile("file")
		var filename = strconv.Itoa(int(time.Now().UnixMilli())) + "_" + h.Filename
		persistedFile, _ := os.Create(path.Join(config.BinDir, filename))
		io.Copy(persistedFile, incomingFile)
		persistedFile.Close()
	}

	retrieveImage := func(w http.ResponseWriter, request *http.Request) {
		files, _ := ioutil.ReadDir(config.BinDir)
		fmt.Println(request.Header)
		incomingLen, err := strconv.Atoi(request.Header["Content-Length"][0])
		fmt.Println(err)
		var buf = make([]byte, incomingLen)
		request.Body.Read(buf)
		fmt.Println(string(buf))
		for _, f := range files {
			fmt.Println(f.Name())
		}
	}

	http.Handle("/", http.FileServer(http.Dir(config.WebpageDir)))
	http.HandleFunc("/upload", saveImage)
	http.HandleFunc("/retrieve", retrieveImage)
	err := http.ListenAndServe(fmt.Sprintf("%s:80", getIP()), nil)
	fmt.Println(err)

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
}

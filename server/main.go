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

	saveImage := func (w http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(32000000)
		incomingFile, h, _ := request.FormFile("file")
		persistedFile, _ := os.Create(path.Join(config.BinDir, h.Filename))
		io.Copy(persistedFile, incomingFile)
		persistedFile.Close()
	}

	http.Handle("/", http.FileServer(http.Dir(config.WebpageDir)))
	http.HandleFunc("/upload", saveImage)
	err := http.ListenAndServe(fmt.Sprintf("%s:80", getIP()), nil)
	fmt.Println(err)

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
}

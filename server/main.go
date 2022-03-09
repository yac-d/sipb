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
)

func saveImage(w http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32000000)
	incomingFile, h, _ := request.FormFile("file")
	fmt.Println(h)
	fmt.Println(incomingFile)
	persistedFile, _ := os.Create(path.Join(".", h.Filename))
	io.Copy(persistedFile, incomingFile)
	persistedFile.Close()
}

func getIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return net.IP(conn.LocalAddr().(*net.UDPAddr).IP).String()
}

func main() {
	wd, _ := os.Getwd()
	http.Handle("/", http.FileServer(http.Dir(path.Join(wd, "webpages"))))
	http.HandleFunc("/upload", saveImage)
	fmt.Println(getIP())
	err := http.ListenAndServe(fmt.Sprintf("%s:80", getIP()), nil)
	fmt.Println(err)

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
}

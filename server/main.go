package main

import (
	"net"
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	//"io"
)

func saveImage(w http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Body)
}

func getIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return net.IP(conn.LocalAddr().(*net.UDPAddr).IP).String()
}

func main() {
	wd, _ := os.Getwd()
	http.Handle("/", http.FileServer(http.Dir(wd)))
	http.HandleFunc("/upload", saveImage)
	fmt.Println(getIP())
	err := http.ListenAndServe(fmt.Sprintf("%s:80", getIP()), nil)
	fmt.Println(err)

	var terminator = make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)
	<-terminator
}

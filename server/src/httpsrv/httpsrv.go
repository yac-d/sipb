package httpsrv

import (
	"net/http"
	"log"
	"strconv"
	"fmt"

	"github.com/Eeshaan-rando/sipb/src/filebin"
	"github.com/Eeshaan-rando/sipb/src/configdef"
)

type HTTPSrv struct {
	filebin filebin.FileBin
	config  configdef.Config
}

func New(conf configdef.Config, fb filebin.FileBin) *HTTPSrv {
	return &HTTPSrv{config: conf, filebin: fb}
}

func (srv *HTTPSrv) handleSave(w http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(64000000)
	incomingFile, h, err := request.FormFile("file")

	log.Printf("Receiving file %s", h.Filename)
	if err != nil {
		log.Println("Error reading uploaded file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if trunc := srv.filebin.SaveFile(incomingFile, h); trunc > 0 { // checking if file had to be truncated
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(strconv.FormatInt(trunc, 10)))
	}
}

func (srv *HTTPSrv) handleGetFileDetails(w http.ResponseWriter, request *http.Request) {
	incomingLen, _ := strconv.Atoi(request.Header["Content-Length"][0])
	var buf = make([]byte, incomingLen)
	request.Body.Read(buf)
	
	whichFile, err := strconv.Atoi(string(buf))
	details, err := srv.filebin.DetailsOfNthNewest(whichFile)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid request for file details")
		return
	}

	w.Write(details.AsJSON())
	log.Printf("File %s requested", details.Path)
}

func (srv *HTTPSrv) handleGetFileCount(w http.ResponseWriter, request *http.Request) {
	cnt := srv.filebin.Count()
	w.Write([]byte(strconv.Itoa(cnt)))
	log.Printf("File count requested, currently at %d", cnt)
}

func (srv *HTTPSrv) Start() error {
	http.Handle("/", http.FileServer(http.Dir(srv.config.WebpageDir)))
	http.HandleFunc("/upload", srv.handleSave)
	http.HandleFunc("/retrieve", srv.handleGetFileDetails)
	http.HandleFunc("/retrieve/fileCount", srv.handleGetFileCount)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", srv.config.BindAddr, srv.config.Port), nil)
}

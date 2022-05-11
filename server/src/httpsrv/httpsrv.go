package httpsrv

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin"
)

type HTTPSrv struct {
	filebin            filebin.FileBin
	config             configdef.Config
	OnSave             OnSaveHandler
	OnCountRequested   OnCountHandler
	OnDetailsRequested OnDetailsHandler
}

type OnSaveHandler func(result filebin.SaveFileResult)
type OnDetailsHandler func(result filebin.FileDetailsResult)
type OnCountHandler func(result filebin.FileCountResult)

func New(conf configdef.Config, fb filebin.FileBin) *HTTPSrv {
	return &HTTPSrv{config: conf, filebin: fb}
}

func (srv *HTTPSrv) handleSave(w http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(64000000)
	incomingFile, h, err := request.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := srv.filebin.SaveFile(incomingFile, h)
	result.Error = err
	if result.TruncatedBytes > 0 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(strconv.FormatInt(result.TruncatedBytes, 10)))
	}

	if srv.OnSave != nil {
		srv.OnSave(result)
	}
}

func (srv *HTTPSrv) handleGetFileDetails(w http.ResponseWriter, request *http.Request) {
	incomingLen, _ := strconv.Atoi(request.Header["Content-Length"][0])
	var buf = make([]byte, incomingLen)
	request.Body.Read(buf)

	whichFile, err := strconv.Atoi(string(buf))
	details, result := srv.filebin.DetailsOfNthNewest(whichFile)
	if err != nil {
		result.Error = err
	}

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(details.AsJSON())

	if srv.OnDetailsRequested != nil {
		srv.OnDetailsRequested(result)
	}
}

func (srv *HTTPSrv) handleGetFileCount(w http.ResponseWriter, request *http.Request) {
	result := srv.filebin.Count()
	w.Write([]byte(strconv.Itoa(result.Count)))

	if srv.OnCountRequested != nil {
		srv.OnCountRequested(result)
	}
}

func (srv *HTTPSrv) Start() error {
	http.Handle("/", http.FileServer(http.Dir(srv.config.WebpageDir)))
	http.HandleFunc("/upload", srv.handleSave)
	http.HandleFunc("/retrieve", srv.handleGetFileDetails)
	http.HandleFunc("/retrieve/fileCount", srv.handleGetFileCount)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", srv.config.BindAddr, srv.config.Port), nil)
}

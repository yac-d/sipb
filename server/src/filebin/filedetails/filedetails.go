package filedetails

import (
	"os"
	"net/http"
	"encoding/json"
)

type FileDetails struct {
	Type string
	Path string
	Size int64
}

func New(fileloc, urlpath string) FileDetails {
	f, _ := os.Open(fileloc)
	var fHeader = make([]byte, 512)
	f.Read(fHeader)
	fInfo, _ := f.Stat()
	f.Close()

	var details FileDetails
	details.Type = http.DetectContentType(fHeader)
	details.Path = urlpath
	details.Size = fInfo.Size()

	return details
}

func (d *FileDetails) AsJSON() []byte {
	j, _ := json.Marshal(d)
	return j
}

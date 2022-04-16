package filedetails

import (
	"os"
	"path"
	"net/http"
	"encoding/json"
)

var fslocation string
var urlpath string

type FileDetails struct {
	Type string
	Path string
	Size int64
}

func NewForFile(filename string) FileDetails {
	f, _ := os.Open(path.Join(fslocation, filename))
	var fHeader = make([]byte, 512)
	f.Read(fHeader)
	fInfo, _ := f.Stat()
	f.Close()

	var details FileDetails
	details.Type = http.DetectContentType(fHeader)
	details.Path = path.Join(urlpath, filename)
	details.Size = fInfo.Size()

	return details
}

func (d *FileDetails) AsJSON() []byte {
	j, _ := json.Marshal(d)
	return j
}

func SetFilesystemLocation(fsloc string) {
	fslocation = fsloc
}

func SetURLPath(urlp string) {
	urlpath = urlp
}

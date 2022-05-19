package filedetails

import (
	"encoding/json"
)

type FileDetails struct {
	Type string
	Path string
	Size int64
}

func (d *FileDetails) AsJSON() []byte {
	j, _ := json.Marshal(d)
	return j
}

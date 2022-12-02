package filedetails

import (
	"encoding/json"
)

type FileDetails struct {
	ID       string `json:"id"`
	Name     string `json:"filename"`
	Location string `json:"location"`
	Size     int64  `json:"bytes"`
	Type     string `json:"mimetype"`
}

func (d *FileDetails) AsJSON() []byte {
	j, _ := json.Marshal(d)
	return j
}

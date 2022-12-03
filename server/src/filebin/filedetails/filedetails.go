package filedetails

import (
	"encoding/json"
	"time"
)

type FileDetails struct {
	ID        string    `json:"ID"`
	Name      string    `json:"Name"`
	Size      int64     `json:"Size"`
	Type      string    `json:"Type"`
	Timestamp time.Time `json:"Timestamp"`
	Location  string    `json:"-"`
}

func (d *FileDetails) AsJSON() []byte {
	j, _ := json.Marshal(d)
	return j
}

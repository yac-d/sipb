package filebin

import (
	"mime/multipart"
	"github.com/Eeshaan-rando/filebin/filedetails"
)

type FileBin interface {
	SaveFile(f multipart.File, h multipart.FileHeader) bool
	Count() int
	DetailsOfNthNewest(n int) (filedetails.FileDetails, error)
}

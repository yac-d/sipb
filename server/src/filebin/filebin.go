package filebin

import (
	"mime/multipart"

	"github.com/Eeshaan-rando/sipb/src/filebin/filedetails"
)

type FileBin interface {
	SaveFile(f multipart.File, h *multipart.FileHeader) int64
	DetailsOfNthNewest(n int) (filedetails.FileDetails, error)
	Count() int
}

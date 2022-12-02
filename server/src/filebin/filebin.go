package filebin

import (
	"mime/multipart"

	"github.com/Eeshaan-rando/sipb/src/filebin/filedetails"
)

type FileBin interface {
	SaveFile(f multipart.File, h *multipart.FileHeader) SaveFileResult
	DetailsOfNthNewest(n int) (filedetails.FileDetails, FileDetailsResult)
	Count() FileCountResult
	RemoveOldFiles() error
}

type SaveFileResult struct {
	Error          error
	TruncatedBytes int64
	Filename       string
	Location       string
}

type FileDetailsResult struct {
	Error    error
	Filename string
}

type FileCountResult struct {
	Error error
	Count int
}

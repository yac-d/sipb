package simplefb

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin"
	"github.com/Eeshaan-rando/sipb/src/filebin/filedetails"
	"github.com/Eeshaan-rando/sipb/src/utils"
)

type SimpleFileBin struct {
	config configdef.Config
}

func New(c configdef.Config) *SimpleFileBin {
	var fb = SimpleFileBin{config: c}
	if !utils.FileExists(c.BinDir) {
		os.MkdirAll(c.BinDir, 0755)
	}
	return &fb
}

func (fb *SimpleFileBin) SaveFile(f multipart.File, h *multipart.FileHeader) filebin.SaveFileResult {
	var filename = strconv.Itoa(int(time.Now().UnixMilli())) + "_" + h.Filename
	var result = filebin.SaveFileResult{TruncatedBytes: 0, Error: nil, Filename: filename}
	persistedFile, err := os.Create(path.Join(fb.config.BinDir, filename))
	if err != nil {
		result.Error = err
	}

	if fb.config.MaxFileSize > -1 {
		io.CopyN(persistedFile, f, fb.config.MaxFileSize)
		l := int64(utils.ReaderLen(f))
		if l > fb.config.MaxFileSize {
			result.TruncatedBytes = l - fb.config.MaxFileSize
		}
	} else {
		io.Copy(persistedFile, f)
	}
	persistedFile.Close()

	err = fb.RemoveOldFiles()
	if err != nil {
		result.Error = err
	}

	return result
}

func (fb *SimpleFileBin) RemoveOldFiles() error {
	var err error = nil
	var files []fs.FileInfo
	if fb.config.MaxFileCnt != -1 {
		files, err = ioutil.ReadDir(fb.config.BinDir)
		var i = 0
		for len(files)-i > fb.config.MaxFileCnt {
			os.Remove(path.Join(fb.config.BinDir, files[i].Name()))
			i += 1
		}
	}
	return err
}

func (fb *SimpleFileBin) Count() (result filebin.FileCountResult) {
	files, err := ioutil.ReadDir(fb.config.BinDir)
	result.Count = len(files)
	result.Error = err
	return
}

func (fb *SimpleFileBin) DetailsOfNthNewest(n int) (fd filedetails.FileDetails, result filebin.FileDetailsResult) {
	files, err := ioutil.ReadDir(fb.config.BinDir)
	if err != nil {
		result.Error = err
	}

	if n > len(files) || n < 1 {
		err = errors.New("Invalid request for details")
		if err != nil {
			result.Error = err
		}
		return
	}

	var filename = files[len(files)-n].Name()
	fd = filedetails.New(path.Join(fb.config.BinDir, filename), path.Join(fb.config.BinPath, filename))
	result.Filename = filename
	return
}

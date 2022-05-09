package simplefb

import (
	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/utils"
	"github.com/Eeshaan-rando/sipb/src/filebin/filedetails"
	"os"
	"log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"time"
	"path"
	"errors"
	"strconv"
)

type SimpleFileBin struct {
	config configdef.Config
}

func New(c configdef.Config) *SimpleFileBin {
	var fb = SimpleFileBin{config: c}
	if !utils.FileExists(c.BinDir) {
		os.MkdirAll(c.BinDir, 0755)
		log.Printf("Creating bin directory %s", c.BinDir)
	}
	return &fb
}

func (fb *SimpleFileBin) SaveFile(f multipart.File, h *multipart.FileHeader) (truncatedBy int64) {
	var filename = strconv.Itoa(int(time.Now().UnixMilli())) + "_" + h.Filename
	persistedFile, _ := os.Create(path.Join(fb.config.BinDir, filename))

	if fb.config.MaxFileSize > -1 {
		io.CopyN(persistedFile, f, fb.config.MaxFileSize)
		l := utils.ReaderLen(f)
		if int64(l) > fb.config.MaxFileSize {
			log.Printf("File %s above set size limit, truncating from %d to %d bytes", filename, l, fb.config.MaxFileSize)
			truncatedBy = int64(l) - fb.config.MaxFileSize
		}
	} else {
		truncatedBy = 0
		io.Copy(persistedFile, f)
	}
	persistedFile.Close()
	log.Printf("File %s saved as %s", h.Filename, filename)

	if fb.config.MaxFileCnt != -1 {
		files, _ := ioutil.ReadDir(fb.config.BinDir)
		var i = 0
		for len(files) - i > fb.config.MaxFileCnt {
			os.Remove(path.Join(fb.config.BinDir, files[i].Name()))
			i += 1
			log.Printf("Removed old file %s", files[i].Name())
		}
	}
	return
}

func (fb *SimpleFileBin) Count() int {
	files, _ := ioutil.ReadDir(fb.config.BinDir)
	return len(files)
}

func (fb *SimpleFileBin) DetailsOfNthNewest(n int) (fd filedetails.FileDetails, err error) {
	files, err := ioutil.ReadDir(fb.config.BinDir)

	if n > len(files) || n < 1 {
		err = errors.New("Invalid request for details")
		return
	}
 
	var filename = files[len(files) - n].Name()
	fd = filedetails.New(path.Join(fb.config.BinDir, filename), path.Join(fb.config.BinPath, filename))
	return
}

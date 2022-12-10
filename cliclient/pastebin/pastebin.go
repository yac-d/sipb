package pastebin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	cntPath      = "retrieve/fileCount"
	uploadPath   = "upload"
	downloadPath = "static"
	detailsPath  = "retrieve"
)

type Pastebin struct {
	url string
}

type FileDetails struct {
	ID        string    `json:"ID"`
	Name      string    `json:"Name"`
	Size      int64     `json:"Size"`
	Type      string    `json:"Type"`
	Timestamp time.Time `json:"Timestamp"`
	Note      string    `json:"Note"`
}

func New(pburl string) *Pastebin {
	return &Pastebin{url: pburl}
}

// Count returns the number of files on the pastebin
func (pb Pastebin) Count() (int, error) {
	urlurlpath, _ := url.JoinPath(pb.url, cntPath)
	resp, err := http.Get(urlurlpath)
	if err != nil {
		return 0, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	return strconv.Atoi(string(bytes))
}

// DetailsOfNthNewest returns the details of the Nth most
// recently uploaded file. N starts from 1.
func (pb Pastebin) DetailsOfNthNewest(n int) (FileDetails, error) {
	var details FileDetails

	urlpath, _ := url.JoinPath(pb.url, detailsPath)
	resp, err := http.Post(urlpath, "text/plain", strings.NewReader(strconv.Itoa(n)))
	if err != nil {
		return details, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &details)
	return details, err
}

func (pb Pastebin) Upload(filepath string, note string) error {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	fileWriter, _ := writer.CreateFormFile("file", path.Base(filepath))

	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	fileWriter.Write(contents)

	writer.WriteField("note", note)
	writer.Close()

	urlpath, _ := url.JoinPath(pb.url, uploadPath)
	req, err := http.NewRequest("POST", urlpath, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := new(http.Client)
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response with nonzero status %s", resp.Status)
	}
	return err
}

// DownloadNth downloads the nth most recently uploaded file
func (pb Pastebin) DownloadNth(n int) error {
	details, err := pb.DetailsOfNthNewest(n)
	if err != nil {
		return err
	}

	urlpath, _ := url.JoinPath(pb.url, downloadPath, details.ID)
	resp, err := http.Get(urlpath)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	contents, _ := io.ReadAll(resp.Body)
	return ioutil.WriteFile(details.Name, contents, 0755)
}

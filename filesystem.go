package fbxapi

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type FSTask struct {
	ID             int    `json:"id"`
	Type           string `json:"type"`
	State          string `json:"state"`
	Error          string `json:"error"`
	CreatedTS      int    `json:"created_ts"`
	StartedTS      int    `json:"started_ts"`
	DoneTS         int    `json:"done_ts"`
	Duration       int    `json:"duration"`
	Progress       int    `json:"progress"`
	ETA            int    `json:"eta"`
	From           string `json:"from"`
	To             string `json:"to"`
	NFiles         int    `json:"nfiles"`
	NFilesDone     int    `json:"nfiles_done"`
	TotalBytes     int    `json:"total_bytes"`
	TotalBytesDone int    `json:"total_bytes_done"`
	CurrBytes      int    `json:"curr_bytes"`
	Rate           int    `json:"rate"`
}

type FileInfo struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	MimeType     string `json:"mimetype"`
	Type         string `json:"type"`
	Size         int    `json:"size"`
	Modification int    `json:"modification"`
	Index        int    `json:"index"`
	Link         bool   `json:"link"`
	Target       string `json:"target"`
	Hidden       bool   `json:"hidden"`
	FolderCount  int    `json:"foldercount"`
	FileCount    int    `json:"filecount"`
}

type FileUpload struct {
	ID         int    `json:"id"`
	Size       int    `json:"size"`
	Uploaded   int    `json:"uploaded"`
	Status     string `json:"status"`
	StartDate  int    `json:"start_date"`
	LastUpdate int    `json:"last_update"`
	UploadName string `json:"upload_name"`
	Dirname    string `json:"dirname"`
}

type FileUploadStartAction struct {
	WSRequest
	Size     int    `json:"size"`
	Dirname  string `json:"dirname"`
	Filename string `json:"filename"`
	Force    string `json:"force"`
}

type FileUploadChunkResponse struct {
	TotalLen  int  `json:"total_len"`
	Complete  bool `json:"complete,omitempty"`
	Cancelled bool `json:"cancelled,omitempty"`
}

func encodePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return base64.StdEncoding.EncodeToString([]byte(path))
}

func boolToIntStr(aBool bool) string {
	if aBool {
		return "1"
	}
	return "0"
}

var TasksEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "fs/tasks/",
}

var LsEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "fs/ls/{{.path}}",
}

var InfoEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "fs/info/{{.path}}",
}

var DlEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "dl/{{.path}}",
}

var UlEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "ws/upload",
}

func (c *Client) Ls(path string, onlyFolder, countSubFolder, removeHidden bool) (respFileInfo []FileInfo, err error) {
	defer panicAttack(&err)

	queryParams := url.Values{}
	queryParams.Set("onlyFolder", boolToIntStr(onlyFolder))
	queryParams.Set("countSubFolder", boolToIntStr(countSubFolder))
	queryParams.Set("removeHidden", boolToIntStr(removeHidden))

	params := map[string]string{
		"path": encodePath(path),
	}

	err = c.Query(LsEP).As(params).WithParams(queryParams).Do(&respFileInfo)
	checkErr(err)
	return
}

func (c *Client) Info(path string) (respFileInfo *FileInfo, err error) {
	defer panicAttack(&err)

	params := map[string]string{
		"path": encodePath(path),
	}

	err = c.Query(InfoEP).As(params).Do(&respFileInfo)
	checkErr(err)
	return
}

func (c *Client) Dl(path string) (resp *http.Response, err error) {
	defer panicAttack(&err)

	params := map[string]string{
		"path": encodePath(path),
	}

	resp, err = c.Query(DlEP).As(params).DoRequest()
	checkErr(err)
	return
}

func (c *Client) Upload(path, destDir string) (err error) {
	defer panicAttack(&err)

	conn, err := c.Query(UlEP).WS()
	checkErr(err)
	defer conn.Close()

	f, err := os.Open(path)
	checkErr(err)
	defer f.Close()

	fi, err := f.Stat()
	checkErr(err)

	reqID := int(time.Now().Unix())

	reqUploadStart := &FileUploadStartAction{
		WSRequest: WSRequest{
			Action:    "upload_start",
			RequestID: reqID,
		},
		Size:     int(fi.Size()),
		Dirname:  encodePath(destDir),
		Filename: fi.Name(),
		Force:    "overwrite",
	}

	err = websocket.JSON.Send(conn, reqUploadStart)
	checkErr(err)

	resp := new(WSResponse)

	err = websocket.JSON.Receive(conn, &resp)
	checkErr(err)

	if !resp.Success {
		return errors.New(resp.Msg)
	}

	buf := make([]byte, 512000)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		checkErr(err)

		err = websocket.Message.Send(conn, buf[:n])
		checkErr(err)

		err = websocket.JSON.Receive(conn, &resp)
		checkErr(err)

		if !resp.Success {
			return errors.New(resp.Msg)
		}
	}

	reqUploadFinalize := &WSRequest{
		Action:    "upload_finalize",
		RequestID: reqID,
	}

	err = websocket.JSON.Send(conn, reqUploadFinalize)
	checkErr(err)

	err = websocket.JSON.Receive(conn, &resp)
	checkErr(err)

	if !resp.Success {
		return errors.New(resp.Msg)
	}

	return
}

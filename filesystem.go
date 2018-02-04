package fbxapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

type FileUploadChunkResult struct {
	TotalLen  int  `json:"total_len"`
	Complete  bool `json:"complete,omitempty"`
	Cancelled bool `json:"cancelled,omitempty"`
}

type FileUploadChunkResponse struct {
	WSResponse
	Result FileUploadChunkResult `json:"result,omitempty"`
}

type ShareLink struct {
	Token   string `json:"token,omitempty"`
	Path    string `json:"path,omitempty"`
	Name    string `json:"name,omitempty"`
	Expire  int    `json:"expire"`
	FullURL string `json:"fullurl,omitempty"`
}

func EncodePath(path string) string {
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

var ShareEP = &Endpoint{
	Verb: HTTP_METHOD_POST,
	Url:  "share_link/",
}

func (c *Client) Ls(path string, onlyFolder, countSubFolder, removeHidden bool) (respFileInfo []FileInfo, err error) {
	defer panicAttack(&err)

	queryParams := url.Values{}
	queryParams.Set("onlyFolder", boolToIntStr(onlyFolder))
	queryParams.Set("countSubFolder", boolToIntStr(countSubFolder))
	queryParams.Set("removeHidden", boolToIntStr(removeHidden))

	params := map[string]string{
		"path": EncodePath(path),
	}

	err = c.Query(LsEP).As(params).WithParams(queryParams).Do(&respFileInfo)
	checkErr(err)
	return
}

func (c *Client) Info(path string) (respFileInfo *FileInfo, err error) {
	defer panicAttack(&err)

	params := map[string]string{
		"path": EncodePath(path),
	}

	err = c.Query(InfoEP).As(params).Do(&respFileInfo)
	checkErr(err)
	return
}

func (c *Client) Dl(path string) (resp *http.Response, err error) {
	defer panicAttack(&err)

	params := map[string]string{
		"path": EncodePath(path),
	}

	resp, err = c.Query(DlEP).As(params).DoRequest()
	checkErr(err)
	return
}

func dispatchRecvError(ctx context.Context, entryCh <-chan *WSResponse, errorCh chan<- bool) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-entryCh:
			if !resp.Success {
				errorCh <- true
				return
			}
		}
	}
}

func uploadMsgReceiver(ctx context.Context, conn *websocket.Conn, dispatcher map[string]chan *WSResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var message []byte
			if err := websocket.Message.Receive(conn, &message); err == nil {
				resp := new(WSResponse)
				err = json.Unmarshal(message, &resp)
				checkErr(err)

				if ch, ok := dispatcher[resp.Action]; ok {
					ch <- resp
				}
			}
		}
	}
}

func sendFile(ctx context.Context, cancel context.CancelFunc, dataRecvCh <-chan *WSResponse, conn *websocket.Conn, path string, reqID int) {
	f, err := os.Open(path)
	checkErr(err)

	defer f.Close()

	errorCh := make(chan bool, 1)
	defer close(errorCh)

	go dispatchRecvError(ctx, dataRecvCh, errorCh)

	buf := make([]byte, 512000)

send_loop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-errorCh:
			cancel()
			return
		default:
			n, err := f.Read(buf)
			if err == io.EOF {
				break send_loop
			}
			checkErr(err)

			err = websocket.Message.Send(conn, buf[:n])
			checkErr(err)
		}
	}

	reqUploadFinalize := &WSRequest{
		Action:    "upload_finalize",
		RequestID: reqID,
	}

	err = websocket.JSON.Send(conn, reqUploadFinalize)
	checkErr(err)
}

func (c *Client) Upload(path, destDir string) (err error) {
	defer panicAttack(&err)

	conn, err := c.Query(UlEP).WS()
	checkErr(err)
	defer conn.Close()

	fi, err := os.Stat(path)
	checkErr(err)

	reqID := int(time.Now().Unix())

	reqUploadStart := &FileUploadStartAction{
		WSRequest: WSRequest{
			Action:    "upload_start",
			RequestID: reqID,
		},
		Size:     int(fi.Size()),
		Dirname:  EncodePath(destDir),
		Filename: fi.Name(),
		Force:    "overwrite",
	}

	dispatcher := map[string]chan *WSResponse{
		"upload_start":    make(chan *WSResponse),
		"upload_data":     make(chan *WSResponse),
		"upload_finalize": make(chan *WSResponse),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go uploadMsgReceiver(ctx, conn, dispatcher)

	err = websocket.JSON.Send(conn, reqUploadStart)
	checkErr(err)

	resp := <-dispatcher["upload_start"]
	if !resp.Success {
		return errors.New(resp.Msg)
	}

	go sendFile(ctx, cancel, dispatcher["upload_data"], conn, path, reqID)

	select {
	case <-ctx.Done():
		checkErr(ctx.Err())
	case <-dispatcher["upload_finalize"]:

	}

	return
}

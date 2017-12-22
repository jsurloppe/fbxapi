package fbxapi

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
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
	RequestID int    `json:"request_id,omitempty"`
	Action    string `json:"action"`
	Size      int    `json:"size"`
	Dirname   string `json:"dirname"`
	Filename  string `json:"filename"`
	Force     string `json:"force"`
}

type FileUploadFinalizeAction struct {
	RequestID int    `json:"request_id,omitempty"`
	Action    string `json:"action"`
}

type FileUploadCancelAction struct {
	RequestID int    `json:"request_id,omitempty"`
	Action    string `json:"action"`
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
	Url:  "fs/info/{{.path}}",
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

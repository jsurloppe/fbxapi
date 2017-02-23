package fbxapi

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
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

func encodePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return base64.StdEncoding.EncodeToString([]byte(path))
}

func (c *Client) Tasks() (respFSTask *FSTask, err error) {
	defer panicAttack(&err)

	resp, err := c.request(HTTP_METHOD_GET, "fs/tasks/", nil)
	checkErr(err)

	respFSTask = new(FSTask)
	err = ResultromResponse(resp, respFSTask)
	checkErr(err)

	return
}

func (c *Client) Ls(path string, onlyFolder, countSubFolder, removeHidden bool) (respFileInfo []FileInfo, err error) {
	defer panicAttack(&err)

	strOnlyFolder := strconv.FormatBool(onlyFolder)
	strCountSubFoder := strconv.FormatBool(countSubFolder)
	strRemoveHidden := strconv.FormatBool(removeHidden)

	url := fmt.Sprintf("fs/ls/%s?onlyFolder=%s&countSubFolder=%s&removeHidden=%s",
		encodePath(path), strOnlyFolder, strCountSubFoder, strRemoveHidden)

	resp, err := c.request(HTTP_METHOD_GET, url, nil)
	checkErr(err)

	err = ResultromResponse(resp, &respFileInfo)
	checkErr(err)

	return
}

func (c *Client) Info(path string) (respFileInfo *FileInfo, err error) {
	defer panicAttack(&err)

	url := fmt.Sprintf("fs/info/%s", path)

	resp, err := c.request(HTTP_METHOD_GET, url, nil)
	checkErr(err)

	respFileInfo = new(FileInfo)
	err = ResultromResponse(resp, respFileInfo)
	checkErr(err)

	return
}

func (c *Client) Dl(path string) (reader io.ReadCloser, err error) {
	defer panicAttack(&err)

	url := fmt.Sprintf("dl/%s", encodePath(path))

	resp, err := c.newRequest(HTTP_METHOD_GET, url, nil)
	checkErr(err)

	return resp.Body, nil
}

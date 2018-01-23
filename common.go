package fbxapi

import (
	"encoding/json"
)

const AUTHHEADER = "X-Fbx-App-Auth"
const CTHEADER = "Content-Type"

const HTTP_METHOD_GET = "GET"
const HTTP_METHOD_POST = "POST"
const HTTP_METHOD_PUT = "PUT"
const HTTP_METHOD_DELETE = "DELETE"

const PROTO_HTTP = "http"
const PROTO_HTTPS = "https"
const PROTO_WS = "ws"
const PROTO_WSS = "wss"

type WSRequest struct {
	RequestID int    `json:"request_id,omitempty"`
	Action    string `json:"action"`
}

type WSResponse struct {
	RequestID int             `json:"request_id,omitempty"`
	Action    string          `json:"action"`
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result"`
	ErrorCode string          `json:"error_code"`
	Msg       string          `json:"msg"`
}

type WSNotification struct {
	Action  string          `json:"action"`
	Success bool            `json:"success"`
	Source  string          `json:"source"`
	Event   string          `json:"event"`
	Result  json.RawMessage `json:"result"`
}

func ResultFromResponse(resp *APIResponse, result interface{}) (err error) {
	defer panicAttack(&err)
	err = json.Unmarshal(resp.Result, result)
	checkErr(err)
	return
}

func SelectRequestMethod(updateMethod string, fn func(interface{}) bool, data interface{}) (method string, body []byte, err error) {
	defer panicAttack(&err)

	method = HTTP_METHOD_GET
	if fn(data) {
		method = updateMethod
		body, err = json.Marshal(data)
		checkErr(err)
	}

	return
}

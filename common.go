package fbxapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/mdns"
)

const AUTHHEADER = "X-Fbx-App-Auth"

const HTTP_METHOD_GET = "GET"
const HTTP_METHOD_POST = "POST"
const HTTP_METHOD_PUT = "PUT"
const HTTP_METHOD_DELETE = "DELETE"

const PROTO_HTTP = "http"
const PROTO_HTTPS = "https"
const PROTO_WS = "ws"

type Response struct {
	Success   bool            `json:"success"`
	Msg       string          `json:"msg"`
	UID       string          `json:"uid"`
	ErrorCode string          `json:"error_code"`
	Result    json.RawMessage `json:"result"`
}

type Client struct {
	http         *http.Client
	mutex        sync.Mutex
	Freebox      *Freebox
	AppID        string
	Version      int
	SessionToken string
}

type Freebox struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	APIVersion
	RespAuthorize
}

func (fb *Freebox) fromServiceEntry(service *mdns.ServiceEntry) {
	fb.DeviceName = strings.TrimSuffix(service.Name, "._fbx-api._tcp.local.")
	fb.DeviceName = strings.Replace(fb.DeviceName, "\\", "", -1)
	fb.Host = service.Host
	fb.Port = service.Port

	for _, field := range service.InfoFields {
		r := strings.Split(field, "=")
		switch r[0] {
		case "api_version":
			fb.APIVersion.APIVersion = r[1]
		case "api_base_url":
			fb.APIBaseURL = r[1]
		case "device_type":
			fb.DeviceType = r[1]
		case "uid":
			fb.UID = r[1]
		case "https_available":
			fb.RemoteHTTPSAvailable = r[1] == "0"
		case "https_port":
			studip, err := strconv.Atoi(r[1])
			checkErr(err)
			fb.RemoteHTTPSPort = studip
			checkErr(err)
		case "api_domain":
			fb.RemoteAPIDomain = r[1]
		}
	}
}

func ResultFromResponse(resp *Response, result interface{}) (err error) {
	defer panicAttack(&err)
	err = json.Unmarshal(resp.Result, &result)
	checkErr(err)
	return
}

// NewClient create a client from a freebox instance
func NewClient(appId string, freebox *Freebox) (client *Client, err error) {
	defer panicAttack(&err)

	iVersion, err := APIVersionToInt(freebox.APIVersion.APIVersion)
	checkErr(err)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	client = &Client{
		Freebox: freebox,
		http:    httpClient,
		Version: iVersion,
		AppID:   appId,
	}

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

func (c *Client) makeUrl(endpoint, proto string) string {
	return fmt.Sprintf("%s://%s:%d/api/v%d/%s", proto, c.Freebox.Host, c.Freebox.Port, c.Version, endpoint)
}

func (c *Client) newRequest(method, endpoint string, body []byte) (resp *http.Response, err error) {
	bodyBuffer := bytes.NewBuffer(body)
	url := c.makeUrl(endpoint, PROTO_HTTPS)

	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	if len(c.SessionToken) > 0 {
		req.Header.Add(AUTHHEADER, c.SessionToken)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	resp, err = client.Do(req)
	return
}

func (c *Client) httpRequest(method, endpoint string, body []byte, openSession bool) (result *Response, err error) {
	defer panicAttack(&err)

	if openSession && len(c.Freebox.AppToken) > 0 && len(c.SessionToken) == 0 {
		c.OpenSession(c.AppID, c.Freebox.AppToken)
	}

	resp, err := c.newRequest(method, endpoint, body)
	checkErr(err)
	if resp.StatusCode == 400 {
		panic(errors.New("Bad request"))
	}

	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	drebug(string(bodyResp))

	result = new(Response)
	err = json.Unmarshal(bodyResp, &result)
	checkErr(err)

	if !result.Success {
		panicStr := fmt.Sprintf("[%s] %s", result.ErrorCode, result.Msg)
		panic(errors.New(panicStr))
	}
	return
}

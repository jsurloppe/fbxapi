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

type Response struct {
	Success   bool            `json:"success"`
	Msg       string          `json:"msg"`
	UID       string          `json:"uid"`
	ErrorCode string          `json:"error_code"`
	Result    json.RawMessage `json:"result"`
}

type Client struct {
	Host         string
	Port         int
	Version      int
	SSL          bool
	SessionToken string
	http         *http.Client
	mutex        sync.Mutex
}

type Freebox struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	APIVersion
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

func ResultromResponse(resp *Response, result interface{}) (err error) {
	defer panicAttack(&err)
	err = json.Unmarshal(resp.Result, &result)
	checkErr(err)
	return
}

func NewClient(host string, port, version int, ssl bool) (client *Client, err error) {
	defer panicAttack(&err)

	if strings.HasSuffix(host, ".") {
		addr, err := MdnsResolve(host)
		checkErr(err)
		host = addr.String()
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: tr}

	client = &Client{
		Host:    host,
		Port:    port,
		Version: version,
		SSL:     ssl,
		http:    httpClient,
	}
	return
}

func NewClientFromFreebox(freebox Freebox, ssl bool) (client *Client, err error) {
	defer panicAttack(&err)

	iVersion, err := APIVersionToInt(freebox.APIVersion.APIVersion)
	checkErr(err)

	if ssl {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		httpClient := &http.Client{Transport: tr}

		client = &Client{
			Host:    freebox.RemoteAPIDomain,
			Port:    freebox.RemoteHTTPSPort,
			Version: iVersion,
			SSL:     true,
			http:    httpClient,
		}
	} else {
		host := freebox.Host
		if strings.HasSuffix(host, ".") {
			addr, err := MdnsResolve(host)
			checkErr(err)
			host = addr.String()
		}

		client = &Client{
			Host:    host,
			Port:    freebox.Port,
			Version: iVersion,
			SSL:     false,
			http:    new(http.Client),
		}
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

func (c *Client) makeUrl(endpoint string) string {
	proto := PROTO_HTTP
	if c.SSL {
		proto = PROTO_HTTPS
	}
	return fmt.Sprintf("%s://%s:%d/api/v%d/%s", proto, c.Host, c.Port, c.Version, endpoint)
}

func (c *Client) newRequest(method, endpoint string, body []byte) (resp *http.Response, err error) {
	bodyBuffer := bytes.NewBuffer(body)
	url := c.makeUrl(endpoint)

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

func (c *Client) request(method, endpoint string, body []byte) (result *Response, err error) {
	defer panicAttack(&err)

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

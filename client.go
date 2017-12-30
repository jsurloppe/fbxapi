package fbxapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"text/template"

	"golang.org/x/net/websocket"
)

type App struct {
	AppID      string
	AppVersion string
}

type Client struct {
	http         *http.Client
	mutex        sync.Mutex
	Freebox      *Freebox
	Version      int
	SessionToken string
	App          *App
}

type APIResponse struct {
	Success   bool            `json:"success"`
	Msg       string          `json:"msg"`
	UID       string          `json:"uid"`
	ErrorCode string          `json:"error_code"`
	Result    json.RawMessage `json:"result"`
}

type Endpoint struct {
	Verb         string
	Url          string
	NoAuth       bool
	BodyRequired bool
	RespStruct   interface{}
}

type Query struct {
	Client         *Client
	Endpoint       *Endpoint
	urlParams      map[string]string
	queryParams    url.Values
	body           []byte
	rawAPIResponse *APIResponse
}

var tmpl *template.Template

func init() {
	tmpl = template.New("url")
}

func NewClient(app *App, fb *Freebox) *Client {
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: tr}

	return &Client{
		Freebox: fb,
		http:    httpClient,
		Version: 4,
		App:     app,
	}
}

func (c *Client) Query(ep *Endpoint) Query {
	return Query{
		Client:   c,
		Endpoint: ep,
	}
}

func (c *Client) WithSession(token string) *Client {
	c.SessionToken = token
	return c
}

func (q Query) As(params map[string]string) Query {
	q.urlParams = params
	return q
}

func (q Query) WithParams(params url.Values) Query {
	q.queryParams = params
	return q
}

func (q Query) WithBody(body interface{}) Query {
	bodyJSON, err := json.Marshal(body)
	checkErr(err)
	q.body = bodyJSON
	return q
}

func (q Query) Inspect(resp *APIResponse) Query {
	q.rawAPIResponse = resp
	return q
}

func (q Query) DoRequest() (resp *http.Response, err error) {
	defer panicAttack(&err)

	if !q.Endpoint.NoAuth && q.Client.SessionToken == "" {
		q.Client.OpenSession(q.Client.App.AppID, q.Client.Freebox.AppToken)
	}
	url := q.makeUrl(PROTO_HTTPS, q.urlParams)
	url.RawQuery = q.queryParams.Encode()

	bodyBuffer := bytes.NewBuffer(q.body)
	req, err := http.NewRequest(q.Endpoint.Verb, url.String(), bodyBuffer)
	checkErr(err)

	if len(q.Client.SessionToken) > 0 {
		req.Header.Add(AUTHHEADER, q.Client.SessionToken)
	}

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	resp, err = client.Do(req)
	checkErr(err)

	err = checkHTTPError(resp)
	checkErr(err)

	return resp, err
}

func (q Query) WS() (conn *websocket.Conn, err error) {
	defer panicAttack(&err)

	if !q.Endpoint.NoAuth && q.Client.SessionToken == "" {
		err = q.Client.OpenSession(q.Client.App.AppID, q.Client.Freebox.AppToken)
		checkErr(err)
	}
	url := q.makeUrl(PROTO_WSS, q.urlParams)
	url.RawQuery = q.queryParams.Encode()

	hostname, err := os.Hostname()
	checkErr(err)

	config, err := websocket.NewConfig(url.String(), "http://"+hostname)
	config.Header = http.Header{}
	config.Header.Set(AUTHHEADER, q.Client.SessionToken)
	config.TlsConfig = tlsConfig

	conn, err = websocket.DialConfig(config)
	checkErr(err)

	return
}

func checkHTTPError(resp *http.Response) error {
	if resp.StatusCode >= 500 {
		return errors.New(resp.Status)
	}
	return nil
}

func checkAPIError(resp *APIResponse) error {
	if !resp.Success {
		return errors.New(resp.Msg)
	}
	return nil
}

func (q Query) Do(endStruct interface{}) (err error) {
	defer panicAttack(&err)
	resp, err := q.DoRequest()
	checkErr(err)

	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	err = json.Unmarshal(bodyResp, &q.rawAPIResponse)
	checkErr(err)

	err = checkAPIError(q.rawAPIResponse)
	checkErr(err)

	if endStruct != nil {
		err = ResultFromResponse(q.rawAPIResponse, endStruct)
		checkErr(err)
	}

	return
}

func (q *Query) makeUrl(proto string, urlmap map[string]string) *url.URL {
	ep := q.Endpoint.Url
	buf := new(bytes.Buffer)
	if urlmap != nil {
		ptmpl, err := tmpl.Parse(q.Endpoint.Url)
		checkErr(err)
		err = ptmpl.Execute(buf, urlmap)
		checkErr(err)
		ep = buf.String()
	}
	return &url.URL{
		Scheme: proto,
		Host:   fmt.Sprintf("%s:%d", q.Client.Freebox.Host, q.Client.Freebox.Port),
		Path:   fmt.Sprintf("/api/v%d/%s", q.Client.Version, ep),
	}
}

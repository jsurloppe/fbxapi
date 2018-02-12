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
	"reflect"
	"strconv"
	"sync"
	"text/template"

	"golang.org/x/net/websocket"
)

type App struct {
	ID      string
	Name    string
	Version string
	Token   string
}

type Client struct {
	http    *http.Client
	mutex   sync.Mutex
	session *Session
}

type Session struct {
	*APIVersion
	*RespSession
	Version int
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
	contentType    string
}

var tmpl *template.Template

func init() {
	tmpl = template.New("url")
}

func (c *Client) Query(ep *Endpoint) Query {
	return Query{
		Client:   c,
		Endpoint: ep,
	}
}

func (c *Client) WithSession(session *Session) *Client {
	c.session = session
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

func (q Query) WithFormBody(body interface{}) Query {
	bodyJSON, err := json.Marshal(body)
	checkErr(err)

	m := make(map[string]interface{})
	err = json.Unmarshal(bodyJSON, &m)
	checkErr(err)

	values := stringify(m)
	q.body = []byte(values.Encode())
	q.contentType = "application/x-www-form-urlencoded"
	return q
}

func (q Query) Inspect(resp *APIResponse) Query {
	q.rawAPIResponse = resp
	return q
}

func (q Query) DoRequest() (resp *http.Response, err error) {
	defer panicAttack(&err)

	url := q.makeUrl(PROTO_HTTPS, q.urlParams)
	url.RawQuery = q.queryParams.Encode()

	bodyBuffer := bytes.NewBuffer(q.body)
	req, err := http.NewRequest(q.Endpoint.Verb, url.String(), bodyBuffer)
	checkErr(err)

	if len(q.Client.session.Token) > 0 {
		req.Header.Add(AUTHHEADER, q.Client.session.Token)
	}

	if q.contentType != "" {
		req.Header.Add(CTHEADER, q.contentType)
	}

	resp, err = q.Client.http.Do(req)
	checkErr(err)

	err = checkHTTPError(resp)
	checkErr(err)

	return resp, err
}

func (q Query) WS() (conn *websocket.Conn, err error) {
	defer panicAttack(&err)

	url := q.makeUrl(PROTO_WSS, q.urlParams)
	url.RawQuery = q.queryParams.Encode()

	hostname, err := os.Hostname()
	checkErr(err)

	config, err := websocket.NewConfig(url.String(), "http://"+hostname)
	config.Header = http.Header{}
	config.Header.Set(AUTHHEADER, q.Client.session.Token)
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
		Host:   fmt.Sprintf("%s:%d", q.Client.session.RemoteAPIDomain, q.Client.session.RemoteHTTPSPort),
		Path:   fmt.Sprintf("%sv%d/%s", q.Client.session.APIBaseURL, q.Client.session.Version, ep),
	}
}

func stringify(in map[string]interface{}) (values url.Values) {
	values = url.Values{}
	for k, v := range in {
		rv := reflect.ValueOf(v)
		var vs string
		switch rv.Interface().(type) {
		case int, int8, int16, int32, int64:
			vs = strconv.FormatInt(rv.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			vs = strconv.FormatUint(rv.Uint(), 10)
		case float32:
			vs = strconv.FormatFloat(rv.Float(), 'f', 4, 32)
		case float64:
			vs = strconv.FormatFloat(rv.Float(), 'f', 4, 64)
		case []byte:
			vs = string(rv.Bytes())
		case string:
			vs = rv.String()
		case bool:
			vs = boolToIntStr(rv.Bool())
		}
		values.Add(k, vs)
	}
	return
}

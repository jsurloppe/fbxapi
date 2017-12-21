package fbxapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"text/template"
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
	Verb              string
	Url               string
	NoAuth            bool
	UrlParamsRequired bool
	BodyRequired      bool
	RespStruct        interface{}
}

type Query struct {
	Client         *Client
	Endpoint       *Endpoint
	urlParams      map[string]string
	body           []byte
	rawAPIResponse *APIResponse
}

var tmpl *template.Template

func init() {
	tmpl = template.New("url")
}

func NewClient(app *App, fb *Freebox) *Client {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
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

func (q Query) With(body interface{}) Query {
	bodyJSON, err := json.Marshal(body)
	checkErr(err)
	q.body = bodyJSON
	return q
}

func (q Query) Inspect(resp *APIResponse) Query {
	q.rawAPIResponse = resp
	return q
}

func (q Query) Do(endStruct interface{}) (err error) {
	defer panicAttack(&err)

	if !q.Endpoint.NoAuth && q.Client.SessionToken == "" {
		q.Client.OpenSession(q.Client.App.AppID, q.Client.Freebox.AppToken)
	}

	url := q.makeUrl(PROTO_HTTPS, q.urlParams)

	bodyBuffer := bytes.NewBuffer(q.body)
	req, err := http.NewRequest(q.Endpoint.Verb, url, bodyBuffer)
	checkErr(err)

	if len(q.Client.SessionToken) > 0 {
		req.Header.Add(AUTHHEADER, q.Client.SessionToken)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	checkErr(err)

	if resp.StatusCode >= 400 {
		panic(errors.New(resp.Status))
	}

	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	err = json.Unmarshal(bodyResp, &q.rawAPIResponse)
	checkErr(err)

	if endStruct != nil {
		err = ResultFromResponse(q.rawAPIResponse, endStruct)
		checkErr(err)
	}

	return
}

func (q *Query) makeUrl(proto string, urlmap map[string]string) string {
	url := q.Endpoint.Url
	buf := new(bytes.Buffer)
	if urlmap != nil {
		ptmpl, err := tmpl.Parse(q.Endpoint.Url)
		checkErr(err)
		err = ptmpl.Execute(buf, urlmap)
		checkErr(err)
		url = buf.String()
	}
	return fmt.Sprintf("%s://%s:%d/api/v%d/%s", proto, q.Client.Freebox.Host, q.Client.Freebox.Port, q.Client.Version, url)
}

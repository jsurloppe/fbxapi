package fbxapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/mdns"
)

type Freebox struct {
	Host string
	Port int
	APIVersion
	Authorization
}

func NewFreebox(host string, port int) *Freebox {
	return &Freebox{
		Host: host,
		Port: port,
	}
}

func (fb *Freebox) NewClient() *Client {
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: tr}

	return &Client{
		http: httpClient,
	}
}

func (app *App) createTokenReq() (tr *TokenRequest, err error) {
	hostname, err := os.Hostname()
	checkErr(err)

	tr = new(TokenRequest)
	tr = &TokenRequest{
		AppId:      app.ID,
		AppName:    app.Name,
		AppVersion: app.Version,
		DeviceName: hostname,
	}
	return
}

func (fb *Freebox) getAPIVersionURL(proto string) *url.URL {
	return &url.URL{
		Scheme: proto,
		Host:   fmt.Sprintf("%s:%d", fb.Host, fb.Port),
		Path:   fmt.Sprintf("/api_version"),
	}
}

func (fb *Freebox) NewSession() (sess *Session, err error) {
	defer panicAttack(&err)

	url := fb.getAPIVersionURL(PROTO_HTTPS)
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(url.String())
	checkErr(err)

	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	version := new(APIVersion)
	err = json.Unmarshal(bodyResp, &version)
	checkErr(err)

	iVersion, err := APIVersionToInt(version.APIVersion)
	checkErr(err)

	sess = &Session{
		APIVersion:  version,
		RespSession: &RespSession{},
		Version:     iVersion,
	}
	return
}

func (fb *Freebox) OpenSession(app *App) (client *Client, err error) {
	defer panicAttack(&err)

	if app.Token == "" {
		checkErr(errors.New("AppToken required"))
	}
	client = fb.NewClient()
	session, err := fb.NewSession()
	checkErr(err)
	respLogin, err := client.WithSession(session).Login()
	checkErr(err)
	password := ComputePassword(app.Token, respLogin.Challenge)
	reqSession := ReqSession{AppId: app.ID, Password: password}
	client.session.RespSession, err = client.Session(reqSession)
	checkErr(err)
	return
}

func (fb *Freebox) Register(app *App) (respAuth *Authorization, err error) {
	defer panicAttack(&err)
	client := fb.NewClient()

	req, err := app.createTokenReq()
	checkErr(err)

	respAuth = new(Authorization)
	err = client.Query(AuthorizeEP).WithBody(req).Do(respAuth)
	checkErr(err)
	return
}

func NewFromServiceEntry(service *mdns.ServiceEntry) (fb *Freebox) {
	fb = new(Freebox)
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
	return
}

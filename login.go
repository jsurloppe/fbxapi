package fbxapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const AUTH_STATUS_GRANTED = "granted"
const AUTH_STATUS_PENDING = "pending"

// ------- Request -----------

type TokenRequest struct {
	AppId      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type ReqSession struct {
	AppId    string `json:"app_id"`
	Password string `json:"password"`
}

func ComputePassword(key, challenge string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(challenge))
	payload := mac.Sum(nil)
	return hex.EncodeToString(payload)
}

// ------- Response ---------

type APIVersion struct {
	UID                  string `json:"uid"`
	DeviceName           string `json:"device_name"`
	DeviceType           string `json:"device_type"`
	APIBaseURL           string `json:"api_base_url"`
	APIVersion           string `json:"api_version"`
	RemoteHTTPSAvailable bool   `json:"https_available"`
	RemoteHTTPSPort      int    `json:"https_port"`
	RemoteAPIDomain      string `json:"api_domain"`
}

type RespAuthorize struct {
	AppToken string `json:"app_token"`
	TrackID  int    `json:"track_id"`
}

type RespAuthorizeTrack struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge"`
}

func (resp *RespAuthorizeTrack) isGranted() bool {
	return resp.Status == AUTH_STATUS_GRANTED
}

func (resp *RespAuthorizeTrack) isPending() bool {
	return resp.Status == AUTH_STATUS_PENDING
}

type RespLogin struct {
	LoggedIn     bool   `json:"logged_in"`
	Challenge    string `json:"challenge"`
	PasswordSalt string `json:"password_salt"`
}

type RespSession struct {
	SessionToken string          `json:"session_token"`
	Challenge    string          `json:"challenge"`
	Permissions  map[string]bool `json:"permissions"`
}

func (c *Client) Authorize(tokenReq TokenRequest) (respAuth *RespAuthorize, err error) {
	defer panicAttack(&err)
	tokenReqJSON, err := json.Marshal(tokenReq)
	checkErr(err)
	resp, err := c.request(HTTP_METHOD_POST, "login/authorize/", tokenReqJSON)
	checkErr(err)
	respAuth = new(RespAuthorize)
	err = ResultFromResponse(resp, respAuth)
	checkErr(err)

	return
}

func (c *Client) TrackLogin(track_id int) (respAuth *RespAuthorizeTrack, err error) {
	defer panicAttack(&err)

	url := fmt.Sprintf("login/authorize/%d", track_id)
	resp, err := c.request(HTTP_METHOD_GET, url, nil)
	checkErr(err)
	respAuth = new(RespAuthorizeTrack)
	err = ResultFromResponse(resp, respAuth)
	checkErr(err)

	return
}

func (c *Client) Login() (respLogin *RespLogin, err error) {
	defer panicAttack(&err)

	resp, err := c.request(HTTP_METHOD_GET, "login/", nil)
	checkErr(err)
	respLogin = new(RespLogin)
	err = ResultFromResponse(resp, respLogin)
	checkErr(err)

	return
}

func (c *Client) Session(reqSess ReqSession) (respSess *RespSession, err error) {
	defer panicAttack(&err)

	reqSessJson, err := json.Marshal(reqSess)
	checkErr(err)
	resp, err := c.request(HTTP_METHOD_POST, "login/session/", reqSessJson)
	checkErr(err)
	respSess = new(RespSession)
	err = ResultFromResponse(resp, respSess)
	checkErr(err)

	return
}

func (c *Client) OpenSession(appID, appToken string) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer panicAttack(&err)

	err = c.Logout()
	checkErr(err)
	respLogin, err := c.Login()
	checkErr(err)
	password := ComputePassword(appToken, respLogin.Challenge)
	reqSession := ReqSession{AppId: appID, Password: password}
	respSession, err := c.Session(reqSession)
	checkErr(err)
	c.SessionToken = respSession.SessionToken
	return
}

func (c *Client) Logout() (err error) {
	defer panicAttack(&err)

	if len(c.SessionToken) > 0 {
		_, err = c.request(HTTP_METHOD_POST, "login/logout/", nil)
		checkErr(err)
		c.SessionToken = ""
	}

	return err
}

package fbxapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
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

type Authorization struct {
	AppToken string `json:"app_token"`
	TrackID  int    `json:"track_id"`
}

type AuthorizationState struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge"`
}

func (resp *AuthorizationState) isGranted() bool {
	return resp.Status == AUTH_STATUS_GRANTED
}

func (resp *AuthorizationState) isPending() bool {
	return resp.Status == AUTH_STATUS_PENDING
}

type RespLogin struct {
	LoggedIn     bool   `json:"logged_in"`
	Challenge    string `json:"challenge"`
	PasswordSalt string `json:"password_salt"`
}

type RespSession struct {
	Token       string          `json:"session_token"`
	Challenge   string          `json:"challenge"`
	Permissions map[string]bool `json:"permissions"`
}

var AuthorizeEP = &Endpoint{
	Verb:         HTTP_METHOD_POST,
	Url:          "login/authorize",
	BodyRequired: true,
}

var TrackAuthorizeEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "login/authorize/{{.track_id}}",
}

func (c *Client) Register(tokenReq *TokenRequest) (respAuth *Authorization, err error) {
	respAuth = new(Authorization)
	err = c.Query(AuthorizeEP).WithBody(tokenReq).Do(respAuth)
	checkErr(err)
	return
}

/*func (c *Client) Authorize(tokenReq TokenRequest) (respAuth *Authorization, err error) {
	defer panicAttack(&err)
	tokenReqJSON, err := json.Marshal(tokenReq)
	checkErr(err)
	resp, err := c.httpRequest(HTTP_METHOD_POST, "login/authorize/", tokenReqJSON, false)
	checkErr(err)
	respAuth = new(Authorization)
	err = ResultFromResponse(resp, respAuth)
	checkErr(err)

	return
}*/

var LoginEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "login/",
}

func (c *Client) Login() (respLogin *RespLogin, err error) {
	defer panicAttack(&err)

	respLogin = new(RespLogin)
	err = c.Query(LoginEP).Do(respLogin)
	checkErr(err)

	return
}

var SessionEP = &Endpoint{
	Verb: HTTP_METHOD_POST,
	Url:  "login/session/",
}

func (c *Client) Session(reqSess ReqSession) (respSess *RespSession, err error) {
	defer panicAttack(&err)

	respSess = new(RespSession)
	err = c.Query(SessionEP).WithBody(reqSess).Do(respSess)
	checkErr(err)

	return
}

var LogoutEP = &Endpoint{
	Verb: HTTP_METHOD_POST,
	Url:  "login/logout/",
}

func (c *Client) Logout() (err error) {
	defer panicAttack(&err)

	if len(c.session.Token) > 0 {
		err = c.Query(LogoutEP).Do(nil)
		checkErr(err)
		c.session.Token = ""
	}

	return err
}

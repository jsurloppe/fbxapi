package fbxapi

type FTPConfig struct {
	Enabled             bool   `json:"enabled"`
	AllowAnonymous      bool   `json:"allow_anonymous"`
	AllowAnonymousWrite bool   `json:"allow_anonymous_write"`
	Password            string `json:"password"`
}

func (c *Client) FTPConfig(ftpConf *FTPConfig) (respFTPConf *FTPConfig, err error) {
	defer panicAttack(&err)

	method, body, err := SelectRequestMethod(HTTP_METHOD_PUT, dataIsNil, ftpConf)
	checkErr(err)
	resp, err := c.request(method, "ftp/config/", body)
	checkErr(err)
	respFTPConf = new(FTPConfig)
	ResultromResponse(resp, respFTPConf)
	return
}

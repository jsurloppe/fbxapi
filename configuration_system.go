package fbxapi

type SystemConfig struct {
	FirmwareVersion  string `json:"firmware_version"`
	Mac              string `json:"mac"`
	Serial           string `json:"serial"`
	Uptime           string `json:"uptime"`
	UptimeVal        int    `json:"uptime_val"`
	BoardName        string `json:"board_name"`
	TempCPUm         int    `json:"temp_cpum"`
	TempSW           int    `json:"temp_sw"`
	TempCPUb         int    `json:"temp_cpub"`
	FanRPM           int    `json:"fan_rpm"`
	BoxAuthenticated bool   `json:"box_authenticated"`
}

func (c *Client) System() (sysConf *SystemConfig, err error) {
	defer panicAttack(&err)
	resp, err := c.request(HTTP_METHOD_GET, "system/", nil)
	checkErr(err)
	sysConf = new(SystemConfig)
	err = ResultromResponse(resp, &sysConf)
	checkErr(err)
	return
}

func (c *Client) Reboot() (err error) {
	defer panicAttack(&err)
	_, err = c.request(HTTP_METHOD_POST, "reboot/", nil)
	checkErr(err)
	return
}

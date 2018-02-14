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
	DiskStatus       string `json:"disk_status"`
	BoxFlavor        string `json:"box_flavor"`
	UserMainStorage  string `json:"user_main_storage"`
}

var SystemEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "system/",
}

var RebootEP = &Endpoint{
	Verb: HTTP_METHOD_POST,
	Url:  "system/reboot/",
}

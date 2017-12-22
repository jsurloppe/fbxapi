package fbxapi

import (
	"strconv"
	"strings"

	"github.com/hashicorp/mdns"
)

type Freebox struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	APIVersion
	Authorization
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

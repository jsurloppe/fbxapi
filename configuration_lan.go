package fbxapi

import (
	"net"
)

type ReqHost struct {
	ID          string `json:"id"`
	PrimaryName string `json:"primary_name"`
}

type ReqWoL struct {
	Mac      string `json:"mac"`
	Password string `json:"password"`
}

type LanConfig struct {
	IP          net.IP `json:"ip"`
	Name        string `json:"name"`
	NameDNS     string `json:"name_dns"`
	NameMDNS    string `json:"name_mdns"`
	NameNETBIOS string `json:"name_netbios"`
	Type        string `json:"type"` // maybe mode? need testing
}

type InterfaceStat struct {
	Name      string `json:"name"`
	HostCount int    `json:"host_count"`
}

type LanHostL2Ident struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type LanHostName struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type LanHostL3Connectivity struct {
	Addr              string `json:"addr"`
	Af                string `json:"af"`
	Active            bool   `json:"active"`
	Reachable         bool   `json:"reachable"`
	LastActivity      int    `json:"last_activity"`
	LastTimeReachable int    `json:"last_time_reachable"`
}

type LanHost struct {
	ID                string                  `json:"id"`
	PrimaryName       string                  `json:"primary_name"`
	HostType          string                  `json:"host_type"`
	PrimaryNameManual bool                    `json:"primary_name_manual"`
	L2Ident           LanHostL2Ident          `json:"l2ident"`
	VendorName        string                  `json:"vendor_name"`
	Persistent        bool                    `json:"persistent"`
	Reachable         bool                    `json:"reachable"`
	LastTimeReachable int                     `json:"last_time_reachable"`
	Active            bool                    `json:"active"`
	LastActivity      int                     `json:"last_activity"`
	Names             []LanHostName           `json:"names"`
	L3Connectivities  []LanHostL3Connectivity `json:"l3connectivities"`
}

// LanConfigEP endpoint definition
// Output: LanConfig
var LanConfigEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "lan/config/",
}

// InterfacesEP endpoint definition
// Output: []InterfaceStat
var InterfacesEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "lan/browser/interfaces/",
}

// InterfaceEP endpoint definition
// Output: []LanHost
var InterfaceEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "lan/browser/{{.iface}}/",
}

// InterfaceHostEP endpoint definition
// Output: LanHost
var InterfaceHostEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "lan/browser/{{.iface}}/{{.host_id}}",
}

// WakeOnLanEP endpoint definition
// Output: nil
var WakeOnLanEP = &Endpoint{
	Verb:         HTTP_METHOD_POST,
	Url:          "lan/wol/{{.iface}}/",
	BodyRequired: true,
}

func (lh *LanHost) GetIPv4s() (ips []string) {
	for _, l3 := range lh.L3Connectivities {
		if l3.Active {
			ips = append(ips, l3.Addr)
		}
	}
	return
}

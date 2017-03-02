package fbxapi

import (
	"fmt"
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

func (lh *LanHost) GetIPv4s() (ips []string) {
	for _, l3 := range lh.L3Connectivities {
		if l3.Active {
			ips = append(ips, l3.Addr)
		}
	}
	return
}

func (c *Client) LanConfig(iface, hostID string, reqLC *LanConfig) (respLanConfig *LanConfig, err error) {
	defer panicAttack(&err)

	method, body, err := SelectRequestMethod(HTTP_METHOD_PUT, dataIsNil, reqLC)
	checkErr(err)
	resp, err := c.request(method, "lan/config/", body)
	checkErr(err)
	respLanConfig = new(LanConfig)
	err = ResultFromResponse(resp, respLanConfig)
	checkErr(err)
	return
}

func (c *Client) Interfaces() (ifaceStats []InterfaceStat, err error) {
	defer panicAttack(&err)

	resp, err := c.request(HTTP_METHOD_GET, "lan/browser/interfaces/", nil)
	checkErr(err)
	err = ResultFromResponse(resp, ifaceStats)
	checkErr(err)
	return
}

func (c *Client) Interface(iface string) (lanHosts []LanHost, err error) {
	defer panicAttack(&err)

	url := fmt.Sprintf("lan/browser/%s/", iface)
	resp, err := c.request(HTTP_METHOD_GET, url, nil)
	checkErr(err)
	err = ResultFromResponse(resp, &lanHosts)
	checkErr(err)
	return
}

func (c *Client) InterfaceHost(iface, hostID string, reqHost *ReqHost) (lanHost *LanHost, err error) {
	defer panicAttack(&err)

	method, body, err := SelectRequestMethod(HTTP_METHOD_PUT, dataIsNil, reqHost)
	checkErr(err)
	url := fmt.Sprintf("lan/browser/%s/%s/", iface, hostID)
	resp, err := c.request(method, url, body)
	checkErr(err)
	lanHost = new(LanHost)
	err = ResultFromResponse(resp, lanHost)
	checkErr(err)
	return
}

func (c *Client) WakeOnLan(iface string) (resp *Response, err error) {
	defer panicAttack(&err)
	url := fmt.Sprintf("lan/wol/%s/", iface)
	resp, err = c.request(HTTP_METHOD_GET, url, nil)
	checkErr(err)
	return
}

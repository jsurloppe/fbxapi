package fbxapi

import (
	"testing"
)

func TestInterfaces(t *testing.T) {
	var data []InterfaceStat
	EndpointTester(t, InterfacesEP, &data, &InterfaceStat{}, nil, nil)
}

func TestLanConfig(t *testing.T) {
	EndpointTester(t, LanConfigEP, &LanConfig{}, &LanConfig{}, nil, nil)
}

func TestInterface(t *testing.T) {
	var data []LanHost
	params := map[string]string{
		"iface": "pub",
	}
	EndpointTester(t, InterfaceEP, &data, &LanHost{}, params, nil)
}

func TestInterfaceHost(t *testing.T) {
	params := map[string]string{
		"iface":   "pub",
		"host_id": "ether-ab:cd:ef:12:34:56",
	}
	EndpointTester(t, InterfaceHostEP, &LanHost{}, &LanHost{}, params, nil)
}

func TestWakeOnLan(t *testing.T) {
	t.SkipNow()
	params := map[string]string{
		"iface": "pub",
	}
	EndpointTester(t, InterfaceHostEP, nil, nil, params, &ReqWoL{Mac: "ab:cd:ef:12:34:56"})
}

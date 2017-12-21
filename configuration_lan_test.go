package fbxapi

import (
	"testing"
)

func TestInterfaces(t *testing.T) {
	var data []InterfaceStat
	EndpointTester(t, InterfacesEP, &data, nil, nil)
}

func TestLanConfig(t *testing.T) {
	EndpointTester(t, LanConfigEP, &LanConfig{}, nil, nil)
}

func TestInterface(t *testing.T) {
	var data []LanHost
	params := map[string]string{
		"iface": "pub",
	}
	EndpointTester(t, InterfaceEP, &data, params, nil)
}

func TestInterfaceHost(t *testing.T) {
	var data []LanHost
	params := map[string]string{
		"iface": "pub",
	}
	testClient.Query(InterfaceEP).As(params).Do(&data)

	params = map[string]string{
		"iface":   "pub",
		"host_id": data[0].ID,
	}
	EndpointTester(t, InterfaceHostEP, &LanHost{}, params, nil)
}

func TestWakeOnLan(t *testing.T) {
	t.SkipNow()
	params := map[string]string{
		"iface": "pub",
	}
	EndpointTester(t, InterfaceHostEP, nil, params, &ReqWoL{Mac: "ab:cd:ef:12:34:56"})
}

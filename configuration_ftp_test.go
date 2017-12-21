package fbxapi

import (
	"fmt"
	"testing"
)

func TestFTPConfig(t *testing.T) {
	EndpointTester(t, CurrentFTPConfigEP, &FTPConfig{}, nil, nil)
}

func TestUpdateFTPConfig(t *testing.T) {
	t.SkipNow()
	data := &FTPConfig{}
	EndpointTester(t, CurrentFTPConfigEP, data, nil, nil)
	data.Enabled = false

	EndpointTester(t, UpdateFTPConfigEP, &FTPConfig{}, nil, data)

	data = &FTPConfig{}
	EndpointTester(t, CurrentFTPConfigEP, data, nil, nil)
	fmt.Printf("%#v\n", data)
}

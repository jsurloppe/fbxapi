package fbxapi

import (
	"testing"
)

func TestConnection(t *testing.T) {
	EndpointTester(t, ConnectionEP, &ConnectionStatus{}, nil, nil)
}

func TestConnectionLog(t *testing.T) {
	var data []ConnectionLog
	EndpointTester(t, ConnectionLogEP, &data, nil, nil)
}

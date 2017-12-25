package fbxapi

import (
	"testing"
)

func TestConnection(t *testing.T) {
	EndpointTester(t, ConnectionEP, &ConnectionStatus{}, nil, nil)
}

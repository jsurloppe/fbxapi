package fbxapi

import (
	"testing"
)

func TestSystem(t *testing.T) {
	EndpointTester(t, SystemEP, &SystemConfig{}, nil, nil)
}

func TestReboot(t *testing.T) {
	t.SkipNow()
	err := NewClient(testApp, testFb).Query(RebootEP).Do(nil)
	failOnError(t, err)
}

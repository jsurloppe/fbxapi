package fbxapi

import (
	"os"
	"strconv"
	"testing"
)

func TestTrackState(t *testing.T) {
	params := map[string]string{
		"track_id": strconv.Itoa(testFb.TrackID),
	}
	EndpointTester(t, TrackAuthorizeEP, &AuthorizationState{}, params, nil)
}

func TestRegister(t *testing.T) {
	t.SkipNow()
	hostname, err := os.Hostname()
	failOnError(t, err)

	tokenReq := &TokenRequest{
		AppId:      testClient.App.AppID,
		AppVersion: testClient.App.AppVersion,
		AppName:    "fbxapi",
		DeviceName: hostname,
	}
	_, err = testClient.Register(tokenReq)
	failOnError(t, err)
}

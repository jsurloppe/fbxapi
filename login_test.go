package fbxapi

import (
	"strconv"
	"testing"
)

func TestTrackState(t *testing.T) {
	params := map[string]string{
		"track_id": strconv.Itoa(testFb.TrackID),
	}
	EndpointTester(t, TrackAuthorizeEP, &AuthorizationState{}, &AuthorizationState{}, params, nil)
}

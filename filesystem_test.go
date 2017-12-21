package fbxapi

import (
	"testing"
)

func TestTasks(t *testing.T) {
	var data []FSTask
	EndpointTester(t, TasksEP, &data, nil, nil)
}

package fbxapi

import (
	"testing"
)

func TestTasks(t *testing.T) {
	var data []FSTask
	EndpointTester(t, TasksEP, &data, nil, nil)
}

func TestLs(t *testing.T) {
	folders, err := testClient.Ls("/", false, false, true)
	failOnError(t, err)

	hasHardDrive := false
	for _, info := range folders {
		if info.Name == "Disque dur" {
			hasHardDrive = true
		}
	}
	if !hasHardDrive {
		t.Fail()
	}
}

func TestInfo(t *testing.T) {
	params := map[string]string{
		"path": "Lw==",
	}
	EndpointTester(t, InfoEP, &FileInfo{}, params, nil)
}

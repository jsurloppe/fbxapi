package fbxapi

import (
	"crypto/sha1"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
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

func TestUpload(t *testing.T) {
	err := testClient.Upload("fixtures/lipsum.txt", "/Disque dur/")
	failOnError(t, err)
}

func TestFSDownload(t *testing.T) {
	resp, err := testClient.Dl("/Disque dur/lipsum.txt")
	failOnError(t, err)

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	failOnError(t, err)

	f, err := os.Open("fixtures/lipsum.txt")
	failOnError(t, err)

	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		failOnError(t, err)
	}

	fuck := sha1.Sum(bytes)

	if string(h.Sum(nil)) != string(fuck[:]) {
		t.Fail()
	}
}

func TestShare(t *testing.T) {
	var data ShareLink
	req := &ShareLink{
		Path:   EncodePath("/Disque dur/lipsum.txt"),
		Expire: int(time.Now().Add(time.Minute).Unix()),
	}
	EndpointTester(t, ShareEP, &data, nil, req)
}

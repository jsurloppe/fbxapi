package fbxapi

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	host := flag.String("host", "mafreebox.freebox.fr", "Freebox host")
	port := flag.Int("port", 443, "Freebox HTTPS port")
	token := flag.String("token", "", "App token, will register app if empty")
	trackID := flag.Int("track_id", -1, "App track ID, will register app if empty")

	flag.Parse()

	testApp = &App{
		AppID:      "com.github.jsurloppe.fbxcli",
		AppVersion: "0",
	}

	testFb = &Freebox{
		Host: *host,
		Port: *port,
		Authorization: Authorization{
			AppToken: *token,
			TrackID:  *trackID,
		},
	}

	testClient = NewClient(testApp, testFb)

	code := m.Run()

	testClient.Logout()

	os.Exit(code)
}

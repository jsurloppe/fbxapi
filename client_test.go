package fbxapi

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	host := flag.String("host", "mafreebox.freebox.fr", "Freebox host")
	port := flag.Int("port", 443, "Freebox HTTPS port")
	token := flag.String("token", "", "App token, will register app if empty")
	trackID := flag.Int("track_id", -1, "App track ID, will register app if empty")
	flag.BoolVar(&doTestRegistration, "register", false, "register freebox and exit")

	flag.Parse()

	testApp = &App{
		AppID:      "com.github.jsurloppe.fbxapi",
		AppVersion: "test",
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

	code := 0

	if doTestRegistration {
		hostname, err := os.Hostname()
		checkErr(err)

		tokenReq := &TokenRequest{
			AppId:      testClient.App.AppID,
			AppVersion: testClient.App.AppVersion,
			AppName:    "fbxapi",
			DeviceName: hostname,
		}

		auth, err := testClient.Register(tokenReq)
		checkErr(err)

		fmt.Printf("Touch the right arrow on the freebox display")

		stateCh := make(chan *AuthorizationState, 1)

		go func() {
			params := map[string]string{
				"track_id": strconv.Itoa(auth.TrackID),
			}
			state := new(AuthorizationState)
			for {
				testClient.Query(TrackAuthorizeEP).As(params).Do(&state)

				if !state.isPending() {
					stateCh <- state
					break
				}
				fmt.Printf(".")
				<-time.After(5 * time.Second)
			}
		}()

		select {
		case state := <-stateCh:
			fmt.Printf("\nstatus: %s\n", state.Status)
			if state.isGranted() {
				fmt.Println("run tests with:")
				fmt.Printf("-token: %s -track_id: %d\n", auth.AppToken, auth.TrackID)
			} else {
				code = 1
			}
		case <-time.After(5 * time.Minute):
			fmt.Println("Timeout, try again")
			code = 1
		}
	} else {
		code = m.Run()
		testClient.Logout()
	}

	os.Exit(code)
}

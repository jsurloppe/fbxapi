package fbxapi

import (
	"testing"
	"time"
)

func TestHttpDiscover(t *testing.T) {
	if fb, err := HttpDiscover("mafreebox.freebox.fr", 443); fb == nil || err != nil {
		t.Fail()
	}
}

func TestMdnsDiscover(t *testing.T) {
	fbChan := make(chan *Freebox)
	MdnsDiscover(fbChan)
	select {
	case <-fbChan:
	case <-time.After(time.Second * 1):
		t.Fail()
	}
}

package fbxapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/miekg/dns"
)

const MULTICASTDNSADDR = "224.0.0.251:5353"
const SERVICE = "_fbx-api._tcp"

// domain := "Freebox-Server.local."

func MdnsResolve(domain string) (host net.IP, err error) {
	defer panicAttack(&err)

	udpAddr, err := net.ResolveUDPAddr("udp4", MULTICASTDNSADDR)
	checkErr(err)

	conn, err := net.ListenMulticastUDP("udp4", nil, udpAddr)
	checkErr(err)
	defer conn.Close()

	timeout := time.Now()
	timeout = timeout.Add(2 * time.Second)
	conn.SetDeadline(timeout)

	msg := new(dns.Msg)
	msg.SetQuestion(domain, dns.TypeA)
	wbuf, err := msg.Pack()
	checkErr(err)

	_, err = conn.WriteToUDP(wbuf, udpAddr)
	checkErr(err)

	rbuf := make([]byte, len(domain)+32)
	_, _, err = conn.ReadFromUDP(rbuf)
	checkErr(err)

	err = dns.IsMsg(rbuf)
	checkErr(err)

	ans := new(dns.Msg)
	err = ans.Unpack(rbuf)
	checkErr(err)

	host = ans.Answer[0].(*dns.A).A
	return
}

func MdnsDiscover() (freeboxs []*Freebox) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for service := range entriesCh {
			freebox := new(Freebox)
			freebox.fromServiceEntry(service)
			freeboxs = append(freeboxs, freebox)
		}
	}()

	mdns.Lookup(SERVICE, entriesCh)
	close(entriesCh)
	return freeboxs
}

func HttpDiscover(host string, port int, ssl bool) (freebox *Freebox, err error) {
	defer panicAttack(&err)

	proto := "http"
	if ssl {
		proto = "https"
	}

	url := fmt.Sprintf("%s://%s:%d/api_version", proto, host, port)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	defer resp.Body.Close()
	drebug("[HttpDiscover] %s", body)
	freebox = new(Freebox)
	err = json.Unmarshal(body, &freebox)
	checkErr(err)
	freebox.Host = host
	freebox.Port = port
	return
}
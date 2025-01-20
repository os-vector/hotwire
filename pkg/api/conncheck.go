package api

import (
	"context"
	"fmt"
	"hotwire/pkg/log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var runningmDNS []string
var runningmDNSmu sync.Mutex

func addRunningmDNS(ip string) {
	runningmDNSmu.Lock()
	runningmDNS = append(runningmDNS, ip)
	log.Debug("new runningmDNS list:", runningmDNS)
	runningmDNSmu.Unlock()
}

func findRunningmDNS(ip string) bool {
	runningmDNSmu.Lock()
	defer runningmDNSmu.Unlock()
	for _, str := range runningmDNS {
		if str == ip {
			log.Debug(ip, "already found in mDNS list")
			return true
		}
	}
	return false
}

func removeRunningmDNS(ip string) {
	runningmDNSmu.Lock()
	defer runningmDNSmu.Unlock()
	var newRunnings []string
	for _, str := range runningmDNS {
		if str != ip {
			newRunnings = append(newRunnings, str)
		} else {
			log.Debug("removing", ip, "from mDNS list")
		}
	}
	log.Debug("new mDNS list:", newRunnings)
	runningmDNS = newRunnings
}

// go through all names which are not active
func startmDNS(ip string) {
	log.Debug("running mDNS...")
	addRunningmDNS(ip)
	var packetConnV4 *ipv4.PacketConn
	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		panic(err)
	}

	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		panic(err)
	}

	packetConnV4 = ipv4.NewPacketConn(l4)
	var packetConnV6 *ipv6.PacketConn

	server, err := mdns.Server(packetConnV4, packetConnV6, &mdns.Config{})
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	answer, src, err := server.QueryAddr(ctx, "Vector-E9B3.local")
	fmt.Println("answer:", answer)
	fmt.Println("src:", src)
	fmt.Println("err:", err)
}

func handleConncheck(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Debug("Incoming conncheck: " + ip)
	if findRunningmDNS(ip) {
		go startmDNS(ip)
	}
	startmDNS(r.RemoteAddr)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func InitConncheck() {
	http.HandleFunc("/ok", handleConncheck)
}

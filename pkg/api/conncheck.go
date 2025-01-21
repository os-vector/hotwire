package api

import (
	"context"
	"hotwire/pkg/log"
	"hotwire/pkg/vars"
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

var mDNSexhaustedIPs []string
var exhaustedmu sync.Mutex

func addRunningmDNS(ip string) {
	runningmDNSmu.Lock()
	runningmDNS = append(runningmDNS, ip)
	log.SuperDebug("new runningmDNS list:", runningmDNS)
	runningmDNSmu.Unlock()
}

func findRunningmDNS(ip string) bool {
	runningmDNSmu.Lock()
	defer runningmDNSmu.Unlock()
	for _, str := range runningmDNS {
		if str == ip {
			log.SuperDebug(ip, "already found in mDNS list")
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
	log.SuperDebug("new mDNS list:", newRunnings)
	runningmDNS = newRunnings
}

func addExhaustedmDNS(ip string) {
	exhaustedmu.Lock()
	mDNSexhaustedIPs = append(mDNSexhaustedIPs, ip)
	log.SuperDebug("new exhaustedmDNS list:", mDNSexhaustedIPs)
	exhaustedmu.Unlock()
}

func findExhaustedmDNS(ip string) bool {
	exhaustedmu.Lock()
	defer exhaustedmu.Unlock()
	for _, str := range mDNSexhaustedIPs {
		if str == ip {
			log.SuperDebug(ip, "already found in exhausted mDNS list")
			return true
		}
	}
	return false
}

func removeExhaustedmDNS(ip string) {
	exhaustedmu.Lock()
	defer exhaustedmu.Unlock()
	var newExhausted []string
	for _, str := range mDNSexhaustedIPs {
		if str != ip {
			newExhausted = append(newExhausted, str)
		} else {
			log.Debug("removing", ip, "from exhausted mDNS list")
		}
	}
	log.Debug("new exhausted mDNS list:", newExhausted)
	mDNSexhaustedIPs = newExhausted
}

// go through all names which are not active
// this handles robot IP changes
func startmDNS(ip string) {
	log.Normal("New bot found on network, but IP isn't in list. Figuring out who it is...")
	addRunningmDNS(ip)
	defer removeRunningmDNS(ip)
	for _, name := range vars.GetAllInactiveNames() {
		log.Debug("running mDNS for " + name)
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		_, src, err := server.QueryAddr(ctx, name+".local")
		defer server.Close()
		if err == nil {
			if src.String() == ip {
				rob, _ := vars.GetRobot("", "", name)
				rob.IP = ip
				vars.SaveRobot(rob)
				vars.SetActive(rob.ESN)
				removeExhaustedmDNS(ip)
				log.Debug("TODO: would read jdocs")
				return
				// TODO: fetch jdocs
			}
		}
	}
	log.Error("New robot is making requests to wire-pod, but no saved names match the robot. IP: " + ip)
	addExhaustedmDNS(ip)
}

func handleConncheck(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.SuperDebug("Incoming conncheck: " + ip)
	rob, err := vars.GetRobot(ip, "", "")
	// if we aren't already trying to find this IP on mDNS, and the attempts aren't exhausted for this IP, startmDNS
	if err != nil {
		if !findRunningmDNS(ip) && !findExhaustedmDNS(ip) {
			go startmDNS(ip)
		}
	} else {
		if !rob.Active {
			log.Debug("TODO: would read jdocs")
		}
		vars.SetActive(rob.ESN)
	}
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func InitConncheck() {
	go vars.StartRobotTicker()
	http.HandleFunc("/ok", handleConncheck)
}

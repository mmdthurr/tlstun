package tunnel

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net"
	"slices"
	"strings"
	"sync"

	"github.com/xtaci/smux"
)

type IdToSession struct {
	mu  sync.Mutex
	Its map[string]*smux.Session
	Is  []string
}

func (its *IdToSession) add(k string, s *smux.Session) {
	its.mu.Lock()
	defer its.mu.Unlock()

	its.Its[k] = s
	if !slices.Contains(its.Is, k) {
		its.Is = append(its.Is, k)
	}
}

func (its *IdToSession) del(k string) {
	its.mu.Lock()
	defer its.mu.Unlock()

	delete(its.Its, k)
	_ = slices.DeleteFunc(its.Is, func(i string) bool {
		return k == i
	})
}

var SrvToIdToSession = make(map[string]*IdToSession)

func (s Srv) MainL() {

	log.Printf("inaddr: %s,  %s", s.Tsrvs[0], s.Tsrvs[1])

	cert, err := tls.LoadX509KeyPair(s.Tlscert, s.Tlskey)
	if err != nil {
		log.Printf("mainl: err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := net.Listen("tcp", s.Laddr)
	if err != nil {
		log.Fatalf("mainl: server: listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("mainl: server: accept: %s", err)
		}

		TlsConn := tls.Server(conn, &conf)

		//io sepration goes here
		go func(tlsconn net.Conn, tsrvs []string, pass string, forwardaddr string) {
			inaddr := strings.Split(tlsconn.RemoteAddr().String(), ":")[0]
			log.Printf("inaddr: %s", inaddr)
			if slices.Contains(tsrvs, inaddr) {
				//handle t
				HandleT(tlsconn, pass)
			} else {
				//handle cli
				HandleCli(tlsconn, forwardaddr)
			}

		}(TlsConn, s.Tsrvs, s.Passwd, s.Forwardaddr)
	}

}

func CheckHost(buff []byte) (string, bool) {

	head_split := strings.Split(string(buff), "\r\n")

	for _, h := range head_split {

		if strings.HasPrefix(h, "Host: ") {

			// customize it based on your domain since my domain be something like this kkdfs.usa.choskosh.cfd then [1] would result in usa
			return strings.Split(strings.TrimPrefix(h, "Host: "), ".")[1], true
		}
	}

	return "", false

}

func HandleCli(Conn net.Conn, ForwardAddr string) {

	Buff := make([]byte, 4096)
	rn, err := Conn.Read(Buff)
	if err != nil {
		return
	}

	host, has_host := CheckHost(Buff)
	if has_host {

		ss, ok := SrvToIdToSession[host]
		if ok {
			for {
				rand_session := rand.Intn(len(ss.Is))
				chosen_session := ss.Its[ss.Is[rand_session]]
				new_stream, err := chosen_session.OpenStream()
				if err != nil {
					chosen_session.Close()
					go ss.del(ss.Is[rand_session])
					continue
				}
				new_stream.Write(Buff[:rn])
				go Proxy(Conn, new_stream)
				break
			}
		} else {

			back_stream, err := net.Dial("tcp", ForwardAddr)
			if err == nil {
				back_stream.Write(Buff[:rn])
				go Proxy(Conn, back_stream)
			}

		}
	}

}

func HandleT(conn net.Conn, passwd string) {

	Buff := make([]byte, 1024)
	conn.Read(Buff)

	hello := strings.Split(string(Buff), "_")
	if hello[0] == passwd {
		session, err := smux.Client(conn, nil)
		if err != nil {
			log.Printf("HandleT: Smux: %s", err)
			return
		}

		_, ok := SrvToIdToSession[hello[2]]
		if !ok {
			SrvToIdToSession[hello[2]] = &IdToSession{Its: make(map[string]*smux.Session)}
		}

		go SrvToIdToSession[hello[2]].add(hello[1], session)

	}
}

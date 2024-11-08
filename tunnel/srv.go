package tunnel

import (
	"crypto/tls"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/yamux"
)

var ptol = make(map[string]LandSession)

func handle_lmain(conn net.Conn, passwd string) {
	var nuint uint32 = 0
	authBuff := make([]byte, 1000)
	rsize, _ := conn.Read(authBuff)
	hello_buff := strings.Split(string(authBuff), "_")
	if hello_buff[0] == passwd {
		session, err := yamux.Client(conn, yamux.DefaultConfig())
		if err != nil {
			log.Printf("failed start yamux client: %s", err)
			return
		}
		ptol[hello_buff[1]] = LandSession{S: session}

	} else {
		keys := make([]string, 0, len(ptol))
		for k := range ptol {
			keys = append(keys, k)
		}
		if len(keys) != 0 {

			p := Nextp(keys, &nuint)

			stream, err := ptol[*p].S.Open()
			if err == nil {
				stream.Write(authBuff[:rsize])
				go Proxy(conn, stream)
			} else {
				ptol[*p].S.Close()
				delete(ptol, *p)
			}

		}
	}
}

func start_pdef80(addr string) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("server: 80 listener error: %s", err)
	}
	var nuint uint32 = 0
	for {
		keys := make([]string, 0, len(ptol))
		for k := range ptol {
			keys = append(keys, k)
		}
		if len(keys) != 0 {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("server: def80 accept: %s", err)
			}

			p := Nextp(keys, &nuint)

			stream, err := ptol[*p].S.Open()
			if err == nil {
				go Proxy(conn, stream)
			} else {
				ptol[*p].S.Close()
				delete(ptol, *p)
			}

		}

	}

}

func (s Srv) StartLmain() {
	cert, err := tls.LoadX509KeyPair(s.Tlscert, s.Tlskey)
	if err != nil {
		log.Printf("err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	go start_pdef80(s.Cliaddr)

	listener, err := net.Listen("tcp", s.Laddr)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)

		}

		tlsConn := tls.Server(conn, &conf)
		go handle_lmain(tlsConn, s.Passwd)
	}
}

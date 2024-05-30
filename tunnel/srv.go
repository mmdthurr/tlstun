package tunnel

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/yamux"
)

var ptol = make(map[string]LandSession)

func handle_lmain(conn net.Conn, passwd string) {

	authBuff := make([]byte, 1000)
	conn.Read(authBuff)
	hello_buff := strings.Split(string(authBuff), "_")
	if hello_buff[0] == passwd {
		session, err := yamux.Client(conn, yamux.DefaultConfig())
		if err != nil {
			log.Fatalf("failed start yamux client: %s", err)

		}
		l, ok := ptol[hello_buff[1]]
		if ok {
			l.L.Close()
		}
		listener, err := net.Listen("tcp", "0.0.0.0:"+hello_buff[1])
		if err != nil {
			fmt.Printf("err raised %s", err)
			return
		}

		ptol[hello_buff[1]] = LandSession{L: listener, S: session}

		for {

			outerconn, err := listener.Accept()
			if err != nil {
				log.Printf("server: accept: 11  %s", err)
				listener.Close()
				conn.Close()
				delete(ptol, hello_buff[1])
				break
			}
			stream, err := session.Open()
			if err != nil {
				listener.Close()
				conn.Close()
				delete(ptol, hello_buff[1])
				break

			}

			go Proxy(outerconn, stream)

		}
	}
}

func start_p80() {

	listener, err := net.Listen("tcp", "0.0.0.0:80")
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
				log.Printf("server: 80 accept: %s", err)
			}

			p := Nextp(keys, &nuint)

			stream, err := ptol[*p].S.Open()
			if err == nil {
				go Proxy(conn, stream)
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

	go start_p80()

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

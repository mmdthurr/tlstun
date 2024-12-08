package tunnel

import (
	"crypto/tls"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/yamux"
)

var ptol = make(map[string]map[string]LandSession)

func handle_lmain(conn net.Conn, passwd string) {
	authBuff := make([]byte, 1000)
	_, _ = conn.Read(authBuff)
	hello_buff := strings.Split(string(authBuff), "_")
	if hello_buff[0] == passwd {
		session, err := yamux.Client(conn, yamux.DefaultConfig())
		if err != nil {
			log.Printf("failed start yamux client: %s", err)
			return
		}

		_, ok := ptol[hello_buff[2]]
		if !ok {
			ptol[hello_buff[2]] = make(map[string]LandSession)
		}
		ptol[hello_buff[2]][hello_buff[1]] = LandSession{S: session}

	} else {

		rsp := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 45\r\n\r\nfaghat heyvona ba ghanon jangal okht migiran!"
		_, _ = conn.Write([]byte(rsp))

	}
}

func start_pdef80(addr string) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("server: 80 listener error: %s", err)
	}
	var nuint uint32 = 0

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: def80 accept: %s", err)
		}

		buff := make([]byte, 4096)
		rn, _ := conn.Read(buff)

		spd := strings.Split(string(buff), "\r\n")
	    // log.Printf("%s", spd[10](buff), "\r\n"))
		for i := 0; i < len(spd); i++ {
			if strings.HasPrefix(spd[i], "Host: ") {
				rhost := strings.TrimPrefix(spd[i], "Host: ")

				// customize it based on your domain since my domain be something like this kkdfs.usa.choskosh.cfd then [1] would result in usa
				pk := strings.Split(rhost, ".")[1]

				log.Printf("%s", pk)

				lls, ok := ptol[pk]
				if ok {

					keys := make([]string, 0, len(lls))
					for k := range lls {
						keys = append(keys, k)
					}
					if len(keys) != 0 {

						p := Nextp(keys, &nuint)

						stream, err := lls[*p].S.Open()
						if err == nil {
							stream.Write(buff[:rn])
							go Proxy(conn, stream)
						} else {
							lls[*p].S.Close()
							delete(lls, *p)
						}

					}

				}

				break
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

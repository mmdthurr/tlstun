package tunnel

import (
	"crypto/tls"
	"log"
	"mmd/tlstun/tunnel/ptcp/ptcp"
	"net"
	"slices"
	"strings"

	"github.com/xtaci/smux"
	// "github.com/hashicorp/yamux"
)

var ptol = make(map[string]map[string]LandSession)
var kl = make(map[string][]string)

//func handle_lmain(conn net.Conn, passwd string, matrixaddr string) {
//
//	var nuint uint32 = 0
//
//	authBuff := make([]byte, 4096)
//	rn, _ := conn.Read(authBuff)
//	hello_buff := strings.Split(string(authBuff), "_")
//	if hello_buff[0] == passwd {
//		session, err := smux.Client(conn, nil)
//		if err != nil {
//			log.Printf("failed start yamux client: %s", err)
//			return
//		}
//
//		_, ok := ptol[hello_buff[2]]
//		if !ok {
//			ptol[hello_buff[2]] = make(map[string]LandSession)
//		}
//		ptol[hello_buff[2]][hello_buff[1]] = LandSession{S: session}
//
//	} else {
//		spd := strings.Split(string(authBuff), "\r\n")
//		// log.Printf("%s", spd[10](buff), "\r\n"))
//		has_host := false
//		for i := 0; i < len(spd); i++ {
//			if strings.HasPrefix(spd[i], "Host: ") {
//				rhost := strings.TrimPrefix(spd[i], "Host: ")
//
//				// customize it based on your domain since my domain be something like this kkdfs.usa.choskosh.cfd then [1] would result in usa
//				pk := strings.Split(rhost, ".")[1]
//				log.Printf("%s", pk)
//
//				lls, ok := ptol[pk]
//				if ok {
//					keys := make([]string, 0, len(lls))
//					for k := range lls {
//						keys = append(keys, k)
//					}
//					if len(keys) != 0 {
//						p := Nextp(keys, &nuint)
//						stream, err := lls[*p].S.OpenStream()
//						if err == nil {
//							stream.Write(authBuff[:rn])
//							go Proxy(conn, stream)
//						} else {
//							lls[*p].S.Close()
//							delete(lls, *p)
//						}
//
//					}
//				} else {
//					// pass it to backend matrix
//					backconn, err := net.Dial("tcp", matrixaddr)
//					if err == nil {
//						backconn.Write(authBuff[:rn])
//						go Proxy(conn, backconn)
//					}
//				}
//
//				has_host = true
//				break
//			}
//		}
//
//		if has_host == false {
//			//pass it to backend matrix
//			backconn, err := net.Dial("tcp", matrixaddr)
//			if err == nil {
//				backconn.Write(authBuff[:rn])
//				go Proxy(conn, backconn)
//			}
//		}
//		//rsp := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 45\r\n\r\nfaghat heyvona ba ghanon jangal okht migiran!"
//		//_, _ = conn.Write([]byte(rsp))
//
//	}
//}

func HandleCli(Conn net.Conn, ForwardAddr string) {

	var nuint uint32 = 0
	authBuff := make([]byte, 4096)
	rn, err := Conn.Read(authBuff)
	if err != nil {
		return
	}

	spd := strings.Split(string(authBuff), "\r\n")
	// log.Printf("%s", spd[10](buff), "\r\n"))
	has_host := false
	for i := 0; i < len(spd); i++ {
		if strings.HasPrefix(spd[i], "Host: ") {
			rhost := strings.TrimPrefix(spd[i], "Host: ")

			// customize it based on your domain since my domain be something like this kkdfs.usa.choskosh.cfd then [1] would result in usa
			pk := strings.Split(rhost, ".")[1]
			log.Printf("%s", pk)

			lls, ok := ptol[pk]
			if ok {
				if len(kl[pk]) != 0 {
					p := Nextp(kl[pk], &nuint)
					stream, err := lls[*p].S.OpenStream()
					if err == nil {
						stream.Write(authBuff[:rn])
						go Proxy(Conn, stream)
					} else {
						lls[*p].S.Close()
						delete(lls, *p)
						slices.DeleteFunc(kl[pk], func(s string) bool {
							return s == *p
						})
					}

				}
			} else {
				// pass it to backend matrix
				backconn, err := net.Dial("tcp", ForwardAddr)
				if err == nil {
					backconn.Write(authBuff[:rn])
					go Proxy(Conn, backconn)
				}
			}

			has_host = true
			break
		}
	}

	if has_host == false {
		//pass it to backend matrix
		backconn, err := net.Dial("tcp", ForwardAddr)
		if err == nil {
			backconn.Write(authBuff[:rn])
			go Proxy(Conn, backconn)
		}
	}

}

func (s Srv) LNt() {

	cert, err := tls.LoadX509KeyPair(s.Tlscert, s.Tlskey)
	if err != nil {
		log.Printf("err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := ptcp.Listen("ptcp", s.Tunaddr)
	//	listener, err := net.Listen("tcp", s.Tunaddr)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			continue
		}

		TlsConn := tls.Server(conn, &conf)
		go func(Conn net.Conn) {
			// passwd_port_nl
			Buff := make([]byte, 1024)
			Conn.Read(Buff)
			hello_buff := strings.Split(string(Buff), "_")
			if hello_buff[0] == s.Passwd {
				session, err := smux.Client(conn, nil)
				if err != nil {
					log.Printf("failed start yamux client: %s", err)
					return
				}

				_, ok := ptol[hello_buff[2]]
				if !ok {
					ptol[hello_buff[2]] = make(map[string]LandSession)
				}

				ptol[hello_buff[2]][hello_buff[1]] = LandSession{S: session}

				if slices.Contains(kl[hello_buff[2]], hello_buff[1]) {
					kl[hello_buff[2]] = append(kl[hello_buff[2]], hello_buff[1])
				}
			}
		}(TlsConn)
	}

}

func (s Srv) LCli() {

	cert, err := tls.LoadX509KeyPair(s.Tlscert, s.Tlskey)
	if err != nil {
		log.Printf("err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := net.Listen("tcp", s.Cliaddr)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
		}

		tlsConn := tls.Server(conn, &conf)
		go HandleCli(tlsConn, s.Forwardaddr)
	}
}

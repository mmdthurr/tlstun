package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/hashicorp/yamux"
)

func proxy(conn1, conn2 net.Conn) {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(conn1, conn2)
	}()
	go func() {
		defer wg.Done()
		io.Copy(conn2, conn1)
	}()

	wg.Wait()
}

type Srv struct {
	Laddr  string
	Passwd string
}

func handle_srv_listener(conn net.Conn, passwd string) {

	authBuff := make([]byte, 1000)
	conn.Read(authBuff)
	hello_buff := strings.Split(string(authBuff), "_")
	if hello_buff[0] == passwd {

		session, err := yamux.Client(conn, nil)
		if err != nil {
			log.Fatalf("failed start yamux client: %s", err)

		}

		listener, err := net.Listen("tcp", "0.0.0.0:"+hello_buff[1])
		if err != nil {
			fmt.Printf("err raised %s", err)
			return
		}
		for {
			outerconn, err := listener.Accept()
			if err != nil {
				log.Printf("server: accept: %s", err)

			}
			stream, err := session.Open()
			if err != nil {

				log.Fatalf("failed: %s", err)

			}

			go proxy(outerconn, stream)

		}
	}
}

func (s Srv) strat_listener() {
	cert, err := tls.LoadX509KeyPair("tls.cert", "tls.key")
	if err != nil {
		log.Printf("err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

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
		handle_srv_listener(tlsConn, s.Passwd)
	}
}

type Cli struct {
	RemoteAddr  string
	ExposedPort string
	Passwd      string
	V2p         string
}

func (c Cli) start_cli() {
	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", c.RemoteAddr)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}
	tlsConn := tls.Client(conn, &conf)
	tlsConn.Write([]byte(fmt.Sprintf("%s_%s_", c.Passwd, c.ExposedPort)))

	sesssion, err := yamux.Server(tlsConn, nil)
	if err != nil {
		log.Fatalf("failed: %s", err)

	}
	for {
		stream, err := sesssion.Accept()
		if err != nil {
			log.Fatalf("failed: %s", err)

		}
		destconn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", c.V2p))
		if err != nil {
			log.Fatalf("failed: %s", err)

		}

		go proxy(destconn, stream)
	}

}

func main() {
	mode := flag.String("m", "s", "mode ")
	v2port := flag.String("v2p", "1080", "port")
	eport := flag.String("ep", "5012", "port")
	passwd := flag.String("passwd", "123456", "passwd")
	addr := flag.String("addr", "0.0.0.0:4443", "addr")

	flag.Parse()
	if *mode == "s" {

		// s := Srv{
		// 	Laddr:  "0.0.0.0:8000",
		// 	Passwd: "hello",
		// }

		s := Srv{
			Laddr:  *addr,
			Passwd: *passwd,
		}
		s.strat_listener()

	} else if *mode == "c" {

		// c := Cli{
		// 	RemoteAddr:  "127.0.0.1:8000",
		// 	ExposedPort: "5012",
		// 	V2p:         "1080",
		// 	Passwd:      "hello",
		// }

		c := Cli{
			RemoteAddr:  *addr,
			ExposedPort: *eport,
			V2p:         *v2port,
			Passwd:      *passwd,
		}
		c.start_cli()

	}
}

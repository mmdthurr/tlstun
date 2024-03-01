package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
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
func start_srv(lp, tip, v2p string) {

	cert, err := tls.LoadX509KeyPair("tls.cert", "tls.key")
	if err != nil {
		log.Printf("err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", lp))

	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	for {
		conn, err := listener.Accept()

		if conn.RemoteAddr().String() != tip {
			return
		}

		if err != nil {
			log.Printf("server: accept: %s", err)

		}

		tlsConn := tls.Server(conn, &conf)

		sesssion, err := yamux.Server(tlsConn, nil)
		if err != nil {
			log.Fatalf("failed: %s", err)

		}
		for {
			stream, err := sesssion.Accept()
			if err != nil {
				log.Fatalf("failed: %s", err)

			}
			destconn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", v2p))
			if err != nil {
				log.Fatalf("failed: %s", err)

			}

			go proxy(destconn, stream)
		}

	}
}

func start_cli(raddr, lp string) {
	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", raddr)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}

	tlsConn := tls.Client(conn, &conf)

	session, err := yamux.Client(tlsConn, nil)
	if err != nil {
		log.Fatalf("failed start client: %s", err)

	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", lp))
	if err != nil {
		log.Fatalf("failed: %s", err)

	}
	for {
		outerconn, err := listener.Accept()
		if err != nil {
			fmt.Printf("err %s", err)
		}

		stream, err := session.Open()
		if err != nil {
			log.Fatalf("failed: %s", err)

		}

		go proxy(outerconn, stream)
	}

}

func main() {
	mode := flag.String("m", "s", "mode ")
	lport := flag.String("p", "5000", "port")
	flag.Parse()
	if *mode == "s" {
		start_srv("4433", "107.172.140.38", "1086")

	} else {

		start_cli("87.248.130.83:4433", *lport)
	}
}

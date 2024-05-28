package tunnel

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

func (c Cli) StartCli() {
	
	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", c.RemoteAddr)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}
	tlsConn := tls.Client(conn, &conf)
	tlsConn.Write([]byte(fmt.Sprintf("%s_%s_", c.Passwd, c.ExposePort)))

	sesssion, err := yamux.Server(tlsConn, yamux.DefaultConfig())
	if err != nil {
		log.Fatalf("failed: %s", err)

	}
	for {
		stream, err := sesssion.Accept()
		if err != nil {
			log.Fatalf("failed: %s", err)

		}
		destconn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", c.Bckp))
		if err != nil {
			log.Fatalf("failed: %s", err)

		}
		go Proxy(destconn, stream)
	}

}

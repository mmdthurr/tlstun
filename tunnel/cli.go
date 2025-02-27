package tunnel

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/xtaci/smux"
	// "github.com/hashicorp/yamux"
)

func (c Cli) StartCli() {

	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", c.RemoteAddr)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	tlsConn := tls.Client(conn, &conf)
	tlsConn.Write([]byte(fmt.Sprintf("%s_%s_%s_", c.Passwd, c.ExposePort, c.NodeName)))

	smuxconf := smux.DefaultConfig()
	smuxconf.KeepAliveTimeout = 2 * time.Second
	sesssion, err := smux.Server(tlsConn, smuxconf)
	if err != nil {
		log.Printf("failed: %s", err)
		return

	}

	//	go func(session *smux.Session) {
	//
	//		for {
	//			_, err := sesssion.Ping()
	//			if err != nil {
	//				session.Close()
	//				break
	//			}
	//			time.Sleep(3 * time.Second)
	//		}
	//
	//	}(sesssion)

	for {
		stream, err := sesssion.AcceptStream()
		if err != nil {
			log.Printf("failed: %s", err)
			break
		}
		destconn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", c.Bckp))
		if err != nil {
			log.Printf("failed: %s", err)
			break

		}
		go Proxy(destconn, stream)
	}

}

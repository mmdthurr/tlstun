package tunnel

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/xitongsys/ptcp/ptcp"
	"github.com/xtaci/smux"
	// "github.com/hashicorp/yamux"
)

func (c Cli) StartCli(interf string) {

	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	ptcp.Init(interf)
	conn, err := ptcp.Dial("ptcp", c.RemoteAddr)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	tlsConn := tls.Client(conn, &conf)
	tlsConn.Write([]byte(fmt.Sprintf("%s_%s_%s_", c.Passwd, c.ExposePort, c.NodeName)))

	smuxconf := smux.DefaultConfig()
	smuxconf.KeepAliveTimeout = 5 * time.Second
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

package tunnel

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/xtaci/smux"
)

func (c Cli) StartCli() {

	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", c.RemoteAddr)
	if err != nil {
		log.Printf("conn: dial: failed: %s", err)
		return
	}

	tlsConn := tls.Client(conn, &conf)

	tlsConn.Write(fmt.Appendf([]byte{}, "%s_%s_%s_", c.Passwd, c.ExposePort, c.NodeName))

	smuxconf := smux.DefaultConfig()
	smuxconf.KeepAliveDisabled = true

	sesssion, err := smux.Server(tlsConn, smuxconf)
	if err != nil {
		log.Printf("session: smux: server: %s", err)
		return
	}

	for {
		stream, err := sesssion.AcceptStream()
		if err != nil {
			log.Printf("stream: session: accept: failed: %s", err)
			break
		}
		v2ray_conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", c.Bckp))
		if err != nil {
			log.Printf("v2rayconn: dial: failed: %s", err)
			break

		}
		go Proxy(v2ray_conn, stream)
	}

}

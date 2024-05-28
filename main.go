package main

import (
	"flag"

	"mmd/tlstun/tunnel"
)

func main() {
	mode := flag.String("m", "s", "mode ")
	clibckport := flag.String("bckp", "1080", "port")
	cliexposedport := flag.String("ep", "5012", "port")
	tlscert := flag.String("cert", "tls.cert", "tls certificate")
	tlskey := flag.String("key", "tls.key", "tls key")
	passwd := flag.String("passwd", "123456", "passwd")
	srvaddr := flag.String("addr", "0.0.0.0:4443", "addr")

	flag.Parse()
	if *mode == "s" {

		s := tunnel.Srv{
			Laddr:   *srvaddr,
			Passwd:  *passwd,
			Tlscert: *tlscert,
			Tlskey:  *tlskey,
		}
		s.StartL()

	} else if *mode == "c" {

		c := tunnel.Cli{
			RemoteAddr: *srvaddr,
			ExposePort: *cliexposedport,
			Bckp:       *clibckport,
			Passwd:     *passwd,
		}
		c.StartCli()

	}
}

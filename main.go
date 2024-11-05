package main

import (
	"flag"

	"mmd/tlstun/tunnel"
)

func main() {
	
	tlscert := flag.String("cert", "tls.cert", "tls certificate")
	tlskey := flag.String("key", "tls.key", "tls key")
	passwd := flag.String("passwd", "123456", "passwd")
	srvaddr := flag.String("r", "0.0.0.0:4443", "addr")

	flag.Parse()

	s := tunnel.Srv{
		Laddr:   *srvaddr,
		Passwd:  *passwd,
		Tlscert: *tlscert,
		Tlskey:  *tlskey,
	}
	s.StartLmain()

}

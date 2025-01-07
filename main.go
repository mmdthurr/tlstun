package main

import (
	"flag"
	"strconv"
	"sync"
	"time"

	"mmd/tlstun/tunnel"
)

func main() {

	mode := flag.String("m", "s", "s server, c cli")

	//server
	tlscert := flag.String("cert", "tls.cert", "tls certificate")
	tlskey := flag.String("key", "tls.key", "tls key")
	passwd := flag.String("passwd", "123456", "passwd")
	srvaddr := flag.String("sr", "0.0.0.0:4443", "addr")
	cliaddr := flag.String("clir", "0.0.0.0:80", "cli addr")

	//cli
	raddr := flag.String("r", "127.0.0.1:443", "remote addr")
	stP := flag.Int("port", 5000, "starting port")
	v2P := flag.String("v2p", "1086", "v2ray port")
	connc := flag.Int("c", 10, "amount of connections")
	//passwd := flag.String("passwd", "123456", "tunnel passwd")
	nodename := flag.String("name", "usa", "tunnel node name")

	flag.Parse()
	if *mode == "s" {
		s := tunnel.Srv{
			Laddr:   *srvaddr,
			Cliaddr: *cliaddr,
			Passwd:  *passwd,
			Tlscert: *tlscert,
			Tlskey:  *tlskey,
		}
		s.StartLmain()
	} else if *mode == "c" {
		var wg sync.WaitGroup
		for p := *stP; p < (*stP + *connc); p++ {
			wg.Add(1)
			go func(p int, remoteaddr, passwd, v2port, noden string) {
				for {
					tunnel.Cli{
						NodeName:   noden,
						RemoteAddr: remoteaddr,
						ExposePort: strconv.Itoa(p),
						Passwd:     passwd,
						Bckp:       v2port,
					}.StartCli()
					time.Sleep(1 * time.Second)
				}
			}(p, *raddr, *passwd, *v2P, *nodename)

		}
		wg.Wait()
	}
}

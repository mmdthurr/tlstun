package tunnel

import (
	"io"
	"log"
	"net"
	"sync"
)

func Proxy(to_srvconn, cliconn net.Conn) {

	var wg sync.WaitGroup
	wg.Add(2)

	defer to_srvconn.Close()
	defer cliconn.Close()

	go func() {
		defer wg.Done()
		_, err := io.Copy(to_srvconn, cliconn)
		if err != nil {
			log.Printf("err conn1 conn2 %s \n", err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(cliconn, to_srvconn)
		if err != nil {
			log.Printf("err conn2 conn1 %s \n", err)
		}
	}()

	wg.Wait()
}

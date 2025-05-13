package tunnel

import (
	"io"
	"log"
	"net"
	"sync"
)

func Proxy(conn1, conn2 net.Conn) {

	var wg sync.WaitGroup
	wg.Add(2)

	defer conn1.Close()
	defer conn2.Close()

	go func() {
		defer wg.Done()
		_, err := io.Copy(conn1, conn2)
		if err != nil {
			log.Printf("err conn1 conn2 %s \n", err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(conn2, conn1)
		if err != nil {
			log.Printf("err conn2 conn1 %s \n", err)
		}
	}()

	wg.Wait()
}

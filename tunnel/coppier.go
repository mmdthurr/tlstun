package tunnel

import (
	"io"
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
		io.Copy(conn1, conn2)
	}()
	go func() {
		defer wg.Done()
		io.Copy(conn2, conn1)
	}()

	wg.Wait()
}

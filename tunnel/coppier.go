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

func Proxy2(v2ray, client net.Conn) {

	var wg sync.WaitGroup
	wg.Add(2)

	defer v2ray.Close()
	defer client.Close()

	go func() {
		defer wg.Done()
		io.Copy(v2ray, client)
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(client, v2ray)
		if err != nil {
			rsp := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 45\r\n\r\nfaghat heyvona ba ghanon jangal okht migiran!"
			_, _ = client.Write([]byte(rsp))
		}
	}()

	wg.Wait()
}

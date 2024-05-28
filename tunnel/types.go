package tunnel

import (
	"net"

	"github.com/hashicorp/yamux"
)

type Srv struct {
	Laddr   string
	Passwd  string
	Tlscert string
	Tlskey  string
}

type Cli struct {
	RemoteAddr string
	ExposePort string
	Passwd     string
	Bckp       string
}

type LandSession struct {
	L net.Listener
	S *yamux.Session
}

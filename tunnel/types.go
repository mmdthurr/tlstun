package tunnel

import "github.com/xtaci/smux"

//"github.com/hashicorp/yamux"

type Srv struct {
	Cliaddr    string
	Laddr      string
	Matrixaddr string
	Passwd     string
	Tlscert    string
	Tlskey     string
}

type Cli struct {
	NodeName   string
	RemoteAddr string
	ExposePort string
	Passwd     string
	Bckp       string
}

type LandSession struct {
	S *smux.Session
}

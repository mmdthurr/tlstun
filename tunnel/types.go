package tunnel

type Srv struct {
	Laddr       string
	Tsrvs       []string
	Forwardaddr string
	Passwd      string
	Tlscert     string
	Tlskey      string
}

type Cli struct {
	NodeName   string
	RemoteAddr string
	ExposePort string
	Passwd     string
	Bckp       string
}

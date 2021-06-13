package rewire

// OpenTunnelRequest defines an inbound request to open a tunnel
type OpenTunnelRequest struct {
	AgentSize int `struc:"int16,sizeof=Agent"`
	Agent     string
	HostSize  int `struc:"int8,sizeof=Host"`
	Host      string
	PortSize  int `struc:"int8,sizeof=Port"`
	Port      string
}

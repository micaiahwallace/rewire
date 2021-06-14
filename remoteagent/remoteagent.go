package remoteagent

import (
	"fmt"
	"io"
	"micaiahwallace/rewire"
	"net"

	"github.com/xtaci/smux"

	"github.com/lunixbochs/struc"
)

// Run runs a control client session
func Run(host, port string) error {

	// Create rewire client with remote agent type
	client, err := rewire.InitClient(rewire.AgentConnType, host, port)
	if err != nil {
		return err
	}

	// Get authenticated connection
	conn := client.GetConnection()

	// Receive connection request
	openReq := &rewire.OpenTunnelRequest{}
	struc.Unpack(*conn, openReq)
	destAddr := net.JoinHostPort(openReq.Host, openReq.Port)
	fmt.Println("open tunnel to", openReq.Host, openReq.Port)

	// Filter out unauthorized requests

	// Create forward session
	if err := ForwardConnection(conn, destAddr); err != nil {
		fmt.Println("Unable to create forward stream session", err)
	}

	return nil
}

// ForwardConnection muxes a tcp connection and joins inbound streams to the destination host via new outbound connection per stream
func ForwardConnection(conn *net.Conn, destAddr string) error {

	// Create smux server to accept connections
	muxer, err := smux.Server(*conn, nil)
	if err != nil {
		fmt.Println("Unable to mux stream")
		return err
	}

	// Loop accept new streams
	for {

		// Accept a new stream
		inconn, err := muxer.AcceptStream()
		if err != nil {
			fmt.Println("Cannot accept stream mux", err)
			continue
		}

		// Create the forward in a new go routine
		go func(inconn *smux.Stream) {

			// Dial to destination
			outconn, err := net.Dial("tcp4", destAddr)
			if err != nil {
				fmt.Println("Cannot dial outbound connection")
				return
			}

			// connect inbound and outbound connections
			io.Copy(inconn, outconn)
			io.Copy(outconn, inconn)
		}(inconn)
	}
}

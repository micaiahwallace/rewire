package rewire

import (
	"net"

	"github.com/lunixbochs/struc"
)

// ReceiveServerRequest unpacks request from client
func ReceiveServerRequest(conn *net.Conn, req *Request) error {

	// Unpack bytes into generic request struct
	if err := struc.Unpack(*conn, req); err != nil {
		return err
	}

	return nil
}

// SendServerResponse packs and sends a response to a client
func SendServerResponse(conn *net.Conn, res interface{}) error {

	// Pack bytes onto the connection
	if err := struc.Pack(*conn, res); err != nil {
		return err
	}

	return nil
}

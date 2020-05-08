package rewire

import (
	"errors"
	"log"
	"micaiahwallace/rewire/rwcrypto"
	"net"

	"github.com/lunixbochs/struc"
)

// Server is a rewire transport server
type Server struct {
	Host     string
	Port     string
	listener *net.Listener
}

// InitServer creates a new server
func InitServer(host, port string) (*Server, error) {

	// validate host param
	if host == "" {
		return nil, errors.New("Host must be defined")
	}

	// validate port param
	if port == "" {
		return nil, errors.New("Port must be defined")
	}

	// init rewire lib
	InitLib()

	// create server
	server := &Server{Host: host, Port: port}

	// Setup keystore
	if err := server.SetupKeys(); err != nil {
		return nil, err
	}

	return server, nil
}

// SetupKeys sets up required server keys
func (server *Server) SetupKeys() error {

	// Check if local key exists
	if !rwcrypto.Keys.Contains(Config.LocalKey) {

		// generate new local key
		newKey, err := rwcrypto.GenerateKey(2048)
		if err != nil {
			return err
		}

		// save to the keystore
		rwcrypto.Keys.SavePrivate(Config.LocalKey, newKey)
	}

	return nil
}

// Listen creates tcp listener
func (server *Server) Listen() error {

	// create and start listener
	addr := net.JoinHostPort(server.Host, server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server.listener = &listener

	// Loop accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("connection error:", err.Error())
			continue
		}
		go server.HandleClient(&conn)
	}
}

// ProcessRequest process client requests
func (server *Server) ProcessRequest(req *Request) {

}

// HandleClient process client connections
func (server *Server) HandleClient(conn *net.Conn) {
	for {
		req := &Request{}
		if err := ReceiveRequest(conn, req); err != nil {
			log.Println("Receive request error:", err.Error())
			continue
		}
		server.ProcessRequest(req)
	}
}

// ReceiveRequest unpacks bytes from connection using struc
func ReceiveRequest(conn *net.Conn, req interface{}) error {
	if err := struc.Unpack(*conn, req); err != nil {
		return err
	}
	return nil
}

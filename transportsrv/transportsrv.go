package transportsrv

import (
	"errors"
	"log"
	"micaiahwallace/rewire"
	"net"
)

// Server is a rewire transport server
type Server struct {
	Host     string
	Port     string
	listener *net.Listener
	agents   map[string]*net.Conn
	clients  map[string]*net.Conn
}

// NewServer creates and starts a new transport server
func NewServer(host, port string) (*Server, error) {

	// validate host param
	// if host == "" {
	// 	return nil, errors.New("host must be defined")
	// }

	// validate port param
	if port == "" {
		return nil, errors.New("port must be defined")
	}

	// init rewire lib
	rewire.InitLib()

	// create server
	server := &Server{Host: host, Port: port}

	// Setup keystore
	if keyStoreErr := rewire.SetupServerKeys(); keyStoreErr != nil {
		return nil, keyStoreErr
	}

	// Bind server to port
	if listenErr := server.Listen(); listenErr != nil {
		return nil, listenErr
	}

	return server, nil
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

// HandleClient process client connections
func (server *Server) HandleClient(conn *net.Conn) {
	for {

		// Create a request struct to contain client request
		req := &rewire.Request{}
		if err := rewire.ReceiveServerRequest(conn, req); err != nil {
			log.Println("Receive request error:", err.Error())
			continue
		}

		// Perform logic based on request type
		// res := server.GetRequestResponse(req)

		// Send encrypted response back

	}
}

// ProcessRequest process client requests
func (server *Server) GetRequestResponse(req *rewire.Request) interface{} {
	switch req.Type {
	case rewire.EncryptedReqType:
		// return server.ProcessRequest()
	case rewire.KeyReqType:
		// keyreq := req.(KeyRequest)
		// HandleKeyRequest(keyreq)
	}
	return nil
}

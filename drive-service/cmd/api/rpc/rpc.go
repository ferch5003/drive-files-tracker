package rpc

import (
	"drive-service/internal/platform/driveaccount"
	"google.golang.org/api/drive/v3"
	"log"
	"net"
	"net/rpc"
)

// Server is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type Server struct {
	ServiceAccount driveaccount.ServiceAccount
	DriveService   *drive.Service
}

func (s *Server) Listen(port string) error {
	log.Println("Starting RPC Server on Port:", port)

	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(listen)

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

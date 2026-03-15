package rpc

import (
	"github.com/AntonioMartinezLopez/enginsight/jrpc"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter"
)

type Server struct {
	jrpcServer *jrpc.Server
}

type ServerConfig struct {
	Service counter.Counter
}

func New(config ServerConfig) *Server {
	handlers := jrpc.Handlers{
		CountService: config.Service,
	}
	jrpcServer := jrpc.NewServer(handlers)

	return &Server{
		jrpcServer: jrpcServer,
	}
}

func (s *Server) Handler() *jrpc.Server {
	return s.jrpcServer
}

func (s *Server) Close() {
	s.jrpcServer.Close()
}

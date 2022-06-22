package server

import (
	"context"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/taaaaakahiro/GraphQL-dataloader-MongoDB/pkg/handler"
	"go.uber.org/zap"
)

type Config struct {
	Log *zap.Logger
}

type Server struct {
	Mux     *http.ServeMux
	Handler http.Handler
	server  *http.Server
	handler *handler.Handler
	log     *zap.Logger
}

func NewServer(registry *handler.Handler, cfg *Config) *Server {
	s := &Server{
		Mux:     http.NewServeMux(),
		handler: registry,
	}
	if cfg != nil {
		if log := cfg.Log; log != nil {
			s.log = log
		}
	}

	s.registerHandler()
	return s
}

func (s *Server) registerHandler() {
	// graph ql
	s.Mux.Handle("/gql", playground.Handler("GraphQL playground", "/query"))
	s.Mux.Handle("/query", s.handler.V1.Query())
}

func (s *Server) Serve(listener net.Listener) error {
	server := &http.Server{
		Handler: cors.Default().Handler(s.Mux),
	}
	s.server = server
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

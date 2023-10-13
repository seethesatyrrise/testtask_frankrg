package server

import (
	"context"
	"net/http"
	"testtask_frankrg/internal/utils"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(serverConfig *ServerConfig, router http.Handler) (*Server, error) {
	s := &Server{
		httpServer: &http.Server{
			Addr:    ":" + serverConfig.ServerPort,
			Handler: router,
		},
	}
	return s, nil
}

func (s *Server) Start() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			utils.Logger.Error(err.Error())
		}
	}()

	utils.Logger.Info("server started")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	utils.Logger.Info("server stopped")
	return nil
}

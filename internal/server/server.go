package server

import (
	"context"
	"fmt"
	"github.com/YuriyNazarov/bannersRotator/internal/app"
	"net/http"
)

type Server struct {
	logger app.Logger
	app    Application
	addr   string
	server *http.Server
}

type Application interface {
	GetBanner(slotId, groupId int) (int, error)
	AddBanner(bannerId, slotId int) error
	DeleteBanner(bannerId, slotId int) error
	RegisterClick(bannerId, slotId, groupId int) error
}

func NewServer(logger app.Logger, app Application, addr string) *Server {
	server := &Server{
		logger: logger,
		app:    app,
		addr:   addr,
	}
	mux := NewMux(app)
	server.server = &http.Server{
		Addr:    addr,
		Handler: mux.mux,
	}
	return server
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("err on server shutdown: %s", err))
	}
	return nil
}

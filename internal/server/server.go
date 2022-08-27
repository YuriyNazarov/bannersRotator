package server

import (
	"context"
	"net/http"

	"github.com/YuriyNazarov/bannersRotator/internal/app"
)

type Server struct {
	logger app.Logger
	app    Application
	addr   string
	server *http.Server
}

type Application interface {
	GetBanner(slotID, groupID int) (int, error)
	AddBanner(bannerID, slotID int) error
	DeleteBanner(bannerID, slotID int) error
	RegisterClick(bannerID, slotID, groupID int) error
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
	return s.server.Shutdown(ctx)
}

package rest

import (
	"avito-trainee/internal/config"
	"context"
	"log"
	"net/http"
	"os"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	logger := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	httpServer := &http.Server{
		Addr:              ":" + cfg.HTTPConfig.Port,
		Handler:           handler,
		ReadTimeout:       cfg.HTTPConfig.ReadTimeout,
		WriteTimeout:      cfg.HTTPConfig.WriteTimeout,
		MaxHeaderBytes:    cfg.HTTPConfig.MaxHeaderMegabytes << 20,
		ReadHeaderTimeout: cfg.HTTPConfig.ReadTimeout,
		IdleTimeout:       cfg.HTTPConfig.ReadTimeout,
		ErrorLog:          logger,
	}

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) RunHttp() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

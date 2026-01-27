package server

import (
	"context"
	"econode-cloud/internal/infra/middleware"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPServer struct {
	srv *http.Server
}

func New(port string, app *Container, l *zap.Logger) *HTTPServer {
	// Create server handler
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.AccessLog(l))

	// Init route
	RegisterRoutes(r, app)

	return &HTTPServer{
		srv: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
	}
}

func (s *HTTPServer) Run() error {
	fmt.Printf("server is running in port %s\n", s.srv.Addr)
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

package collector

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spaghettifunk/pixel-collector/collector/middlewares"
	"github.com/spaghettifunk/pixel-collector/collector/routes"
	"github.com/spaghettifunk/pixel-collector/pkg/kafka"
)

// Server encapsulate the web-server object
type Server struct {
	App *echo.Echo
}

// NewServer returns a new instance of the object server
func NewServer(kc *kafka.Client) (*Server, error) {
	// setup echo webserver
	e := echo.New()

	// Middleware
	// inject custom context
	pc, err := middlewares.NewPixelContext(kc)
	if err != nil {
		return nil, err
	}
	e.Use(pc.CustomContext)
	e.Use(middleware.Recover())

	e.GET("/", routes.Version)
	e.GET("/healthz", routes.Healthz)
	e.GET("/collect", routes.Collect)

	return &Server{App: e}, nil
}

// ListenAndServe initialize the server
func (s *Server) ListenAndServe(host, port string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	return s.App.Start(addr)
}

// Shutdown terminates serving the requests
func (s *Server) Shutdown() error {
	return s.App.Shutdown(context.Background())
}

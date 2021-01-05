package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/spaghettifunk/pixel-collector/pkg/kafka"
)

// PixelContext is a custom context that holds the kafka client
type PixelContext struct {
	echo.Context
	kc *kafka.Client
}

// NewPixelContext .
func NewPixelContext(kc *kafka.Client) (*PixelContext, error) {
	return &PixelContext{kc: kc}, nil
}

// CustomContext injects the custom context as middleware
func (l *PixelContext) CustomContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		l.Context = c
		return next(l)
	}
}

// GetKafkaClient returns the internal object of the kafka client
func (l *PixelContext) GetKafkaClient() *kafka.Client {
	return l.kc
}

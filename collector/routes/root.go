package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spaghettifunk/pixel-collector/collector/utils"
)

const (
	// APIVersion contains the current version
	APIVersion = "v0.0.1"
)

// Healthz is a simple endpoint to test if the api is alive
func Healthz(c echo.Context) error {
	r := &utils.Response{Data: "I'm healthy"}
	return c.JSON(http.StatusOK, r)
}

// Version returns the current version of the API
func Version(c echo.Context) error {
	r := &utils.Response{Data: fmt.Sprintf("Pixel-Collector API version: %s", APIVersion)}
	return c.JSON(http.StatusOK, r)
}

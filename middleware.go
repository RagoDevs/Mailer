package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) FilterIPAddress (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.RealIP() != app.config.mail.allowed_ip {
			return echo.NewHTTPError(http.StatusUnauthorized,
				fmt.Sprintf("IP address %s not allowed", c.RealIP()))
		}

		return next(c)
	}
}

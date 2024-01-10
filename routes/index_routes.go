package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckInvoice(c echo.Context) error {
	RefuseExpiredInvoices()
	return c.String(http.StatusOK, "all delete")
}

package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckInvoice(c echo.Context) error {
	if err := RefuseExpiredInvoices(); err != nil {
		return c.String(http.StatusBadRequest, "Error handling the expired invoices")
	}
	return c.String(http.StatusOK, "Refuse expired invoices done correctly")
}

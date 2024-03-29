package routes

import (
	"API_Rest/db"
	"API_Rest/middleware"
	"API_Rest/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func GetInvoicesHandler(c echo.Context) error {
	var invoices []models.Invoice

	if err := db.DataBase().Find(&invoices).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	// There must be a better solution to hide the Issuers pk
	for i := range invoices {
		invoices[i].IssuerPk = ""
	}

	return c.JSON(http.StatusOK, invoices)
}

func BuyInvoiceHandler(c echo.Context) error {
	var params struct {
		InvoiceId     string `validate:"required,uuid"`
		PurchaseFunds int    `validate:"required"`
	}

	// Routine bind and validate params
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	if err := c.Validate(params); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	// Getting the clamins to obtain the user id
	claims, ok := middleware.OptainTokenClaims(c)
	if !ok {
		return c.String(http.StatusInternalServerError, "Error handling tokens claims")
	}

	// Checking that the user exist
	var user models.User
	if err := db.DataBase().Where("user_id = ?", claims["user_id"].(string)).Find(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	tx := db.DataBase().Begin()
	defer tx.Rollback()

	// Getting and checking invoice status
	var invoice models.Invoice
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "invoice_id = ?", params.InvoiceId).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	if invoice.Status != "open" {
		return c.String(http.StatusNotFound, "Invoice not open")
	}
	if !invoice.AllowendPurcharseFunds(params.PurchaseFunds) {
		return c.String(http.StatusForbidden, "Invalid purcharse funds")
	}

	// Getting and checking if the investor can buy the invoice
	var investor models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&investor, "user_id = ?", claims["user_id"]).Error; err != nil {
		return c.String(http.StatusNotFound, "Investor not found")
	}
	if investor.Funds < params.PurchaseFunds {
		return c.String(http.StatusBadRequest, "Insufficient funds")
	}

	// Removing investor funds
	investor.Funds -= params.PurchaseFunds
	if err := tx.Exec("UPDATE users SET funds = ? WHERE user_id = ?", investor.Funds, investor.UserId).Error; err != nil {
		return c.String(http.StatusBadRequest, "Error updating user")
	}

	// Adding funds to the invoice and check if is full
	invoice.Sales(params.PurchaseFunds)
	invoice.Sold()
	if err := tx.Exec("UPDATE invoices SET funds = ?, status = ? WHERE invoice_id = ?", invoice.Funds, invoice.Status, invoice.InvoiceId).Error; err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Creating and validating new purcharse record
	invoiceRecord := models.NewRecord(params.InvoiceId, investor.UserId, params.PurchaseFunds)
	if err := c.Validate(invoiceRecord); err != nil {
		return c.String(http.StatusUnprocessableEntity, "Error validating invoice record")
	}
	if err := tx.Create(invoiceRecord).Error; err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	tx.Commit()

	return c.String(http.StatusOK, "Transaction comitted")

}

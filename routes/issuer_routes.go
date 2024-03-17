package routes

import (
	"API_Rest/db"
	"API_Rest/middleware"
	"API_Rest/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func CreateInvoiceHandler(c echo.Context) error {
	// Getting and validating invoice parametters
	var invoice models.Invoice
	if err := c.Bind(&invoice); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if invoice.Price <= 0 {
		return c.String(http.StatusBadRequest, "Unexpected price")
	}

	claims, ok := middleware.OptainTokenClaims(c)
	if !ok {
		return c.String(http.StatusInternalServerError, "Error reading token claims")
	}
	var user models.User
	if err := db.DataBase().Where("user_id = ?", claims["user_id"].(string)).First(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	invoice.Status = "open"
	invoice.IssuerPk = claims["user_id"].(string)

	if err := c.Validate(invoice); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	if err := db.DataBase().Create(&invoice).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "Created")
}

func ApproveInvoiceHandler(c echo.Context) error {
	param := c.Param("id")

	tx := db.DataBase().Begin()
	defer tx.Rollback()

	// Getting and checking invoice status
	var invoice models.Invoice
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "invoice_id = ?", param).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	if invoice.Status != "waiting" {
		return c.String(http.StatusBadRequest, "Invoice not waiting")
	}

	// Checking that the purchase record and the invoice's total price match
	var invoiceRecordFunds int
	if err := tx.Table("invoice_records").Where("invoice_pk = ?", param).Select("sum(invested_funds)").Scan(&invoiceRecordFunds).Error; err != nil {
		return c.String(http.StatusNotFound, "Error finding ivestred funds register")
	}
	if invoiceRecordFunds != invoice.Funds {
		return c.String(http.StatusBadRequest, "Invoice funds != register funds")
	}

	// Giving funds to the issuer
	var issuerFunds int
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("users").Where("user_id = ?", invoice.IssuerPk).Select("funds").Scan(&issuerFunds).Error; err != nil {
		return c.String(http.StatusNotFound, "Issuer not found")
	}
	issuerFunds += invoiceRecordFunds
	if err := tx.Table("users").Where("user_id = ?", invoice.IssuerPk).Update("funds", issuerFunds).Error; err != nil {
		return c.String(http.StatusBadRequest, "Fail saving issuer data")
	}

	// Closing invoice
	if err := tx.Table("invoices").Where("invoice_id = ?", invoice.InvoiceId).Update("status", "close").Error; err != nil {
		return c.String(http.StatusBadRequest, "Fail saving invoice data")
	}

	tx.Commit()
	return c.String(http.StatusOK, "Invoice aprpoved")

}

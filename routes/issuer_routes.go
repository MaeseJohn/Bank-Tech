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

	gorm := db.DataBase().Table("users").Where("user_id = ?", claims["user_id"].(string)).Select("COUNT(*)")
	if gorm.Error != nil {
		return c.String(http.StatusInternalServerError, gorm.Error.Error())

	}
	if gorm.RowsAffected == 0 {
		return c.String(http.StatusNotFound, "Issuer not found")
	}

	invoice.Status = "open"
	invoice.IssuerPk = claims["user_id"].(string)

	if err := c.Validate(invoice); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	err := db.DataBase().Create(&invoice).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, err.Error())
}

func ApproveInvoiceHandler(c echo.Context) error {
	param := c.Param("id")

	tx := db.DataBase().Begin()
	defer tx.Rollback()

	//Obtaining invoice
	var invoice models.Invoice
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "invoice_id = ?", param).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if invoice.Status != "waiting" {
		return c.String(http.StatusBadRequest, "Invoice not waiting")
	}

	var invoiceRecordFunds int
	if err := tx.Table("invoice_records").Where("invoice_pk = ?", param).Select("sum(invested_funds)").Scan(&invoiceRecordFunds).Error; err != nil {
		return c.String(http.StatusNotFound, "Error finding ivestred funds register")
	}

	if invoiceRecordFunds != invoice.Funds {
		return c.String(http.StatusBadRequest, "Invoice funds != register funds")
	}

	//Giving funds to issuer
	var issuerFunds int
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("users").Where("user_id = ?", invoice.IssuerPk).Select("funds").Scan(&issuerFunds).Error; err != nil {
		return c.String(http.StatusNotFound, "Issuer not found")
	}

	issuerFunds += invoiceRecordFunds

	if err := tx.Table("users").Where("user_id = ?", invoice.IssuerPk).Update("funds", issuerFunds).Error; err != nil {
		return c.String(http.StatusBadRequest, "Fail saving issuer data")
	}

	//Closing invoice

	if err := tx.Table("invoices").Where("invoice_id = ?", invoice.InvoiceId).Update("status", "close").Error; err != nil {
		return c.String(http.StatusBadRequest, "Fail saving invoice data")
	}

	tx.Commit()
	return c.String(http.StatusOK, "Invoice aprpoved")

}

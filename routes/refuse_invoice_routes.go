package routes

import (
	"API_Rest/db"
	"API_Rest/models"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func returnInvestorMoney(record models.InvoiceRecord, tx *gorm.DB) error {
	var userFunds int
	var err error
	if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("users").Where("user_id = ?", record.InvestorPk).Select("funds").Scan(&userFunds).Error; err != nil {
		return err
	}

	userFunds += record.InvestedFunds

	if err = tx.Table("users").Where("user_id = ?", record.InvestorPk).Update("funds", userFunds).Error; err != nil {
		return err
	}

	return err
}

// Search and refuse expired invoices
func RefuseExpiredInvoices() error {
	currentTime := time.Now().UTC().Format(time.DateOnly)
	var invoices []string
	if err := db.DataBase().Table("invoices").Select("invoice_id").Find(&invoices, "expire_date <= ? AND status <> ?", currentTime, "close").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	for _, i := range invoices {

		tx := db.DataBase().Begin()
		defer tx.Rollback()
		var invoice models.Invoice

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invoice, "invoice_id = ? AND status <> ?", i, "close").Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		var invoiceRecord []models.InvoiceRecord
		if err := tx.Where("invoice_pk = ?", invoice.InvoiceId).Find(&invoiceRecord).Error; err != nil {
			return err
		}

		for _, r := range invoiceRecord {
			if err := returnInvestorMoney(r, tx); err != nil {
				return err
			}
		}

		if err := tx.Table("invoices").Where("invoice_id = ?", invoice.InvoiceId).Update("status", "close").Error; err != nil {
			return err
		}

		tx.Commit()
	}
	return nil

}

// Star a ticker with 24h duration
func refuseInvoicesTicker() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		RefuseExpiredInvoices()
	}
}

// Call the ticker at 00:01:00 UTC
func CallTicker() {
	//Ticker functions
	startTime, _ := time.Parse(time.TimeOnly, time.Now().UTC().Format(time.TimeOnly))
	endTime, _ := time.Parse(time.TimeOnly, "23:59:00")
	_ = time.AfterFunc(endTime.Sub(startTime).Abs()+(2*time.Minute), refuseInvoicesTicker)
}

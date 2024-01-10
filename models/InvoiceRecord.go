package models

type InvoiceRecord struct {
	InvoicePk     string `validate:"required,uuid"`
	InvestorPk    string `validate:"required,uuid"`
	InvestedFunds int    `validate:"required"`
}

// Create new *InvoiceRecord object Parameters: InvoicePK, InvestorPK, InvestedFunds
func NewRecord(invoicePk, investorPk string, investedFunds int) InvoiceRecord {
	var invoiceRecord InvoiceRecord
	invoiceRecord.InvoicePk = invoicePk
	invoiceRecord.InvestorPk = investorPk
	invoiceRecord.InvestedFunds = investedFunds

	return invoiceRecord
}

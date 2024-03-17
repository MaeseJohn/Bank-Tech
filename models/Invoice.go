package models

type Invoice struct {
	InvoiceId  string `validate:"required,uuid"`
	IssuerPk   string `validate:"required,uuid"`
	Name       string `validate:"required,alphanum"`
	Price      int    `validate:"required"`
	Funds      int
	Status     string `validate:"oneof=open waiting close"`
	ExpireDate string `validate:"required,datetime=2006-01-02,alloweddate"`
}

// Create invoice request constructor
func NewInvoice(invoiceId, name, expiredate string, price int) *Invoice {
	var invoice Invoice
	invoice.InvoiceId = invoiceId
	invoice.Name = name
	invoice.ExpireDate = expiredate
	invoice.Price = price
	return &invoice
}

// Returns true if the prucharseFunds are valid to buy the invoice
func (i *Invoice) AllowendPurcharseFunds(prucharseFunds int) bool {
	available := i.Price - i.Funds
	return available >= prucharseFunds
}

// If the invoice has been fully purcharsed it status changes to "waiting"
func (i *Invoice) Sold() {
	if i.Price == i.Funds {
		i.Status = "waiting"
	}
}

// Increases the amount of funds indicated by parameter
func (i *Invoice) Sales(x int) {
	i.Funds += x
}

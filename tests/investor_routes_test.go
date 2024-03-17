package tests

import (
	"API_Rest/models"
	"math/rand"
	"net/http"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tjgq/broadcast"
	"gopkg.in/go-playground/assert.v1"
)

func TestBuyInvoice(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers(t)
	status := CreateInvoiceRequest(&validInvoices[0], issuerToken)
	require.Equal(t, http.StatusCreated, status)

	var buyParameters = []struct {
		name      string
		invoiceId string
		funds     int
		token     string
		status    int
	}{

		// Invalid tokens
		{"Zero value token", validInvoices[0].InvoiceId, 101, "", http.StatusUnauthorized},
		{"Invalid token", validInvoices[0].InvoiceId, 102, invalidToken, http.StatusUnauthorized},
		{"Issuer token", validInvoices[0].InvoiceId, 100, issuerToken, http.StatusForbidden},
		{"Unresgister Usertoken", validInvoices[0].InvoiceId, 100, unregisterInvestorToken, http.StatusNotFound},
		// Invalid uuid
		{"Not uuid", "hola", 1000, investorToken, http.StatusUnprocessableEntity},
		{"Unregister uuid", "4ddbb37b-efd5-4564-a2ba-c4ac80925b9f", 100, investorToken, http.StatusNotFound},
		// To many funds
		{"To many funds", validInvoices[0].InvoiceId, 1000, investorToken, http.StatusForbidden},
		// Valid purcharse
		{"Valid purcharse", validInvoices[0].InvoiceId, 100, investorToken, http.StatusOK},
		{"Valid purcharse", validInvoices[0].InvoiceId, 400, investorToken, http.StatusOK},
		// Closed invoice
		{"Closed invoice", validInvoices[0].InvoiceId, 100, investorToken, http.StatusNotFound},
	}

	for _, p := range buyParameters {
		t.Run(p.name, func(t *testing.T) {
			resp := BuyInvoice(p.invoiceId, p.token, p.funds)
			assert.Equal(t, resp, p.status)
		})
	}
}

func TestGetInvoices(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers(t)
	for _, i := range validInvoices {
		status := CreateInvoiceRequest(&i, issuerToken)
		require.Equal(t, http.StatusCreated, status)
	}
	status := BuyInvoice(validInvoices[2].InvoiceId, investorToken, 500)
	require.Equal(t, http.StatusOK, status)

	validInvoices[2].Funds = 500
	invoices, StatusCode := GetInvoicesRequest()

	validInvoices[0].Status = "open"
	validInvoices[1].Status = "open"
	validInvoices[2].Status = "waiting"
	invoices[0].ExpireDate = validInvoices[0].ExpireDate
	invoices[1].ExpireDate = validInvoices[1].ExpireDate
	invoices[2].ExpireDate = validInvoices[2].ExpireDate

	assert.Equal(t, StatusCode, http.StatusOK)
	assert.Equal(t, invoices, validInvoices)
}

func TestConcurrentInvoicePurchase(t *testing.T) {
	CreateEnviroment()
	var randomUsers = []*models.User{
		models.NewUser(uuid.NewString(), "miguel@gmail.com", "miguel", "miguel", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "jose@gmail.com", "jose", "jose", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "pepe@gmail.com", "pepe", "pepe", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "manuel@gmail.com", "manuel", "manuel", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "max@gmail.com", "max", "max", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "josian@gmail.com", "josian", "josian", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "rosa@gmail.com", "rosa", "rosa", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "darius@gmail.com", "darius", "darius", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "gareen@gmail.com", "gareen", "gareen", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "jarvan@gmail.com", "jarvan", "jarvan", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "dimitry@gmail.com", "dimitry", "dimitry", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "vlad@gmail.com", "vlad", "vlad", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "morgan@gmail.com", "morgan", "morgan", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "olivia@gmail.com", "olivia", "olivia", "investor", ((rand.Intn(20) + 5) * 100)),
		models.NewUser(uuid.NewString(), "alicia@gmail.com", "alicia", "alicia", "investor", ((rand.Intn(20) + 5) * 100)),
		//randomUsers[15] == issuers
		models.NewUser(uuid.NewString(), "jorge@gmail.com", "jorge", "jorge", "issuer", 100),
		models.NewUser(uuid.NewString(), "vicente@gmail.com", "vicente", "vicente", "issuer", 100),
		models.NewUser(uuid.NewString(), "marcos@gmail.com", "marcos", "marcos", "issuer", 100),
		models.NewUser(uuid.NewString(), "sara@gmail.com", "sara", "sara", "issuer", 100),
		models.NewUser(uuid.NewString(), "emilia@gmail.com", "emilia", "emilia", "issuer", 100),
		models.NewUser(uuid.NewString(), "pocoyo@gmail.com", "pocoyo", "pocoyo", "issuer", 100),
	}

	var randomInvoices [20]*models.Invoice
	for i := range randomInvoices {
		randomInvoices[i] = models.NewInvoice(uuid.NewString(), "invoice", "2050-01-01", (rand.Intn(50)+10)*100)
	}

	// CREATE USERS AND COUNT TOTAL FUNDS
	var initialFunds int
	for _, u := range randomUsers {
		_ = CreateUserRequest(u)
		initialFunds += u.Funds
	}

	// CREATING INVOICES WITH RANDOM ISSUERS
	for _, i := range randomInvoices {
		randomIssuer := *randomUsers[rand.Intn(6)+15]
		_, tokenString := LoginUserRequest(randomIssuer.Email, randomIssuer.Password)
		_ = CreateInvoiceRequest(i, tokenString)
	}

	// Creating a broadcaster to difetent go routines to call it at the same time and made a concurrency test
	var b broadcast.Broadcaster
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		l := b.Listen()
		invoice := randomInvoices[rand.Intn(16)]
		investor := randomUsers[rand.Intn(15)]
		_, tokenString := LoginUserRequest(investor.Email, investor.Password)
		invested := (rand.Intn(5) + 1) * 100
		wg.Add(1)

		go func() {
			defer wg.Done()
			for range l.Ch {
				BuyInvoice(invoice.InvoiceId, tokenString, invested)
			}
		}()
	}
	b.Send("OPEN")
	b.Close()
	wg.Wait()

	var finalFunds int
	getUsers := GetUsersRequest()
	for _, u := range getUsers {
		finalFunds += u.Funds
	}
	getInvoices, _ := GetInvoicesRequest()
	for _, i := range getInvoices {
		finalFunds += i.Funds
	}
	assert.Equal(t, initialFunds, finalFunds)

}

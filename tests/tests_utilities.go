package tests

import (
	"API_Rest/migrations"
	"API_Rest/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// Tokens
var issuerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjo2MDAwLCJ1c2VyX2lkIjoiYWUyMTU1OTItNWM2NS0xMWVlLThjOTktMDI0MmFjMTIwMDAyIiwidXNlcl9yb2xlIjoiaXNzdWVyIn0.kONVM9itQ_UtCRHkWf60hk3KrhF0LfGpd-3TsHdZAas"
var investorToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjo3MDAwMCwidXNlcl9pZCI6IjY1MjBhMTQ4LTg4ZjUtNDVhNC1hMTg3LTg4M2NmZWJjZjk4NSIsInVzZXJfcm9sZSI6ImludmVzdG9yIn0.Jc7ZBxGu4KXWhbZiXX_n8D1n3eoqLCuurfo6fjdw2n0"
var invalidToken = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.jYW04zLDHfR1v7xdrW3lCGZrMIsVe0vWCfVkN2DRns2c3MN-mcp_-RE6TN9umSBYoNV-mnb31wFf8iun3fB6aDS6m_OXAiURVEKrPFNGlR38JSHUtsFzqTOj-wFrJZN4RwvZnNGSMvK3wzzUriZqmiNLsG8lktlEn6KA4kYVaM61_NpmPHWAjGExWv7cjHYupcjMSmR8uMTwN5UuAwgW6FRstCJEfoxwb0WKiyoaSlDuIiHZJ0cyGhhEmmAPiCwtPAwGeaL1yZMcp0p82cpTQ5Qb-7CtRov3N4DcOHgWYk6LomPR5j5cCkePAz87duqyzSMpCB0mCOuE3CU2VMtGeQ"
var unregisterInvestorToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjo3MDAwMCwidXNlcl9pZCI6IjNlNWYxNmM3LWJhY2EtNGI5Yi05MGExLWQ0MmFlMThlYmU3OCIsInVzZXJfcm9sZSI6ImludmVzdG9yIn0.fyHkxjVsuWAU5RichKRg0eciUMlr2-94435kjoJNPLk"
var unregisterIssuerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjo3MDAwMCwidXNlcl9pZCI6IjAzZDM1OTgzLTI3NGItNGNkMS1hNTNjLWU4Y2RhNjc2Y2IyNyIsInVzZXJfcm9sZSI6Imlzc3VlciJ9.xILjVYIAo2HwEv1FUqEpEUSZv_vBHgouY1BeKij5QaE"

// Valid users for data base
var validUsers = []models.User{
	*models.NewUser("ae215592-5c65-11ee-8c99-0242ac120002", "jonathan@gmail.com", "jonathan", "jonathan", "issuer", 1000),
	*models.NewUser("87224e57-188e-4211-bbb2-c27dca02fa10", "loris@gmail.com", "loris", "loris", "issuer", 600),
	*models.NewUser("6520a148-88f5-45a4-a187-883cfebcf985", "elisa@gmail.com", "elisa", "elisa", "investor", 10000),
}

// Valid invoices
var validInvoices = []models.Invoice{
	*models.NewInvoice(uuid.NewString(), "Invoice1", "2030-01-01", 500),
	*models.NewInvoice(uuid.NewString(), "Invoice1", "2030-01-01", 500),
	*models.NewInvoice(uuid.NewString(), "Invoice1", "2030-01-01", 500),
}

// Expected output of get user
// These are the same users than valid users
// The order of users is the same as that of valid users
var getUsersData = []models.User{
	*models.NewUser("", "jonathan@gmail.com", "", "jonathan", "issuer", 1000),
	*models.NewUser("", "loris@gmail.com", "", "loris", "issuer", 600),
	*models.NewUser("", "elisa@gmail.com", "", "elisa", "investor", 10000),
}

func ResetDB() {
	migrations.CreateDataBase()
	migrations.Migrations.Drop()
	migrations.CreateDataBase()
}

func CreateEnviroment() {
	//Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	ResetDB()
}

func CreateValidUsers(t *testing.T) {
	for _, u := range validUsers {
		status := CreateUserRequest(&u)
		require.Equal(t, http.StatusCreated, status)
	}
}

// Makes a create user request
// Returns response status e.g:201 Created
func CreateUserRequest(user *models.User) int {
	postBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:3000/users/new", "application/json", responseBody)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// Makes a get user request
// Returns respons user
func GetUserRequest(email string) (*models.User, int) {
	//Request
	resp, err := http.Get("http://localhost:3000/users/" + email)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}

	var user models.User

	//Read and close body
	defer resp.Body.Close()
	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Error while reading the response bytes:", err)
		}
		//Unmarshal body to user
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &user, resp.StatusCode
}

// Makes a get users request
// Returns respons users
func GetUsersRequest() []models.User {
	//Request
	resp, err := http.Get("http://localhost:3000/users")
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	//Read and close body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error while reading the response bytes:", err)
	}
	//Unmarshal body to user
	var users []models.User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Fatalln(err)
	}

	return users
}

func LoginUserRequest(email, password string) (status int, token string) {
	var loginParameters struct {
		Email    string
		Password string
	}
	loginParameters.Email = email
	loginParameters.Password = password
	postBody, _ := json.Marshal(&loginParameters)
	responseBody := bytes.NewBuffer(postBody)

	//Request
	resp, err := http.Post("http://localhost:3000/users/login", "application/json", responseBody)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, ""
	}

	//Read and close body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error while reading the response bytes:", err)
	}

	var jwtoken string
	err = json.Unmarshal(body, &jwtoken)
	if err != nil {
		log.Fatalln(err)
	}

	return resp.StatusCode, jwtoken

}

// Create invoice request
// Returns status code
func CreateInvoiceRequest(invoice *models.Invoice, token string) int {
	uri := "http://localhost:3000/invoice"
	postBody, _ := json.Marshal(invoice)
	responseBody := bytes.NewBuffer(postBody)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token

	// Create a new request using http
	req, err := http.NewRequest("POST", uri, responseBody)
	if err != nil {
		log.Fatalln("Error on request.\n[ERROR] -", err)
	}

	// Add content-type header to the req
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func GetInvoicesRequest() ([]models.Invoice, int) {
	// Create a new request using http
	req, err := http.NewRequest("GET", "http://localhost:3000/invoices", nil)
	if err != nil {
		log.Fatalln("Error on request.\n[ERROR] -", err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + investorToken

	// Add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	var invoices []models.Invoice

	//Read body
	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Error while reading the response bytes:", err)
		}
		//Unmarshal body to user
		err = json.Unmarshal(body, &invoices)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return invoices, resp.StatusCode
}

func BuyInvoice(invoiceId, token string, funds int) int {
	var bodyparams struct {
		InvoiceId     string
		PurchaseFunds int
	}
	bodyparams.InvoiceId = invoiceId
	bodyparams.PurchaseFunds = funds

	uri := "http://localhost:3000/invoices"
	postBody, _ := json.Marshal(&bodyparams)
	responseBody := bytes.NewBuffer(postBody)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token

	// Create a new request using http
	req, err := http.NewRequest("POST", uri, responseBody)
	if err != nil {
		log.Fatalln("Error on request.\n[ERROR] -", err)
	}

	// Add content-type header to the req
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode

}

func ApproveInvoceRequest(invoiceId, issuerToken string) int {

	uri := "http://localhost:3000/invoices/" + invoiceId

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + issuerToken

	// Create a new request using http
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		log.Fatalln("Error on request.\n[ERROR] -", err)
	}

	// Add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode

}

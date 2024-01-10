package main

import (
	"API_Rest/db"
	"API_Rest/middleware"
	"API_Rest/migrations"
	"API_Rest/models"
	"API_Rest/routes"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	//Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Fail loading .env file", err)
	}

	migrations.CreateDataBase()

	// Conecting to db
	db.Connection()

	// Creating web server
	e := echo.New()

	//Creating echo validator with go-playgraund/validator library
	v := &models.CustomValidator{}
	e.Validator = v.New(validator.New())

	routes.CallTicker()

	//////////////////////
	//////ENDPOINTS///////
	//////////////////////
	e.GET("/cleanexpiredinvoices", routes.CheckInvoice)

	// User Endpoints
	e.POST("/users/new", routes.CreateUserHandler)
	e.GET("/users", routes.GetUsersHandler)
	e.GET("/users/:email", routes.GetUserHandler)
	e.POST("/users/login", routes.LoginUserHandler)

	//Issuer Endpoints
	e.POST("/invoice", routes.CreateInvoiceHandler, echojwt.JWT([]byte(os.Getenv("SECRET_KEY"))), middleware.RoleValidator("issuer"))
	e.POST("/invoices/:id", routes.ApproveInvoiceHandler, echojwt.JWT([]byte(os.Getenv("SECRET_KEY"))), middleware.RoleValidator("issuer"), middleware.InvoiceOwnerValidator())

	//Investor Endpoints
	e.GET("/invoices", routes.GetInvoicesHandler, echojwt.JWT([]byte(os.Getenv("SECRET_KEY"))), middleware.RoleValidator("investor"))
	e.POST("/invoices", routes.BuyInvoiceHandler, echojwt.JWT([]byte(os.Getenv("SECRET_KEY"))), middleware.RoleValidator("investor"))

	e.Logger.Fatal(e.Start(":3000"))
}

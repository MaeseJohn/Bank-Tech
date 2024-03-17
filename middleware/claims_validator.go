package middleware

import (
	"API_Rest/db"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// False = error occur
func OptainTokenClaims(c echo.Context) (jwt.MapClaims, bool) {
	reqToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	// Check err handling
	token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

// Validates if the user's role is the same as the role passed by parameter
func RoleValidator(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := OptainTokenClaims(c)
			if !ok {
				return c.String(http.StatusInternalServerError, "Fail optaining claims")
			}
			if claims["user_role"] != role {
				return c.String(http.StatusForbidden, "You don't have the permissions to access this resource")
			}
			return next(c)
		}
	}
}

// Validate if the invoice belongs to the user
func InvoiceOwnerValidator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Obtain claims
			claims, ok := OptainTokenClaims(c)
			if !ok {
				return c.String(http.StatusBadRequest, "Fail optaining claims")
			}

			// Obtaining path parameter
			var params struct {
				invoiceId string `validate:"required,uuid"`
			}
			params.invoiceId = c.Param("id")
			if err := c.Validate(params); err != nil {
				return c.String(http.StatusUnprocessableEntity, err.Error())
			}

			var issuer struct {
				IssuerPk string
			}
			err := db.DataBase().Table("invoices").Where("invoice_id = ?", params.invoiceId).Select("issuer_pk").First(&issuer).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return c.String(http.StatusNotFound, "Invoice not found")
				}
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// Check id
			if claims["user_id"] != issuer.IssuerPk {
				return c.String(http.StatusUnauthorized, "You are not the owner of this invoice")
			}

			return next(c)
		}
	}
}

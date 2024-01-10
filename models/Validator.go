package models

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) New(v *validator.Validate) *CustomValidator {
	v.RegisterValidation("alloweddate", allowedDate)
	cv.validator = v
	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Checks that the user's date is later than the current UTC date
func allowedDate(fl validator.FieldLevel) bool {
	userDate, err := time.Parse(time.DateOnly, fl.Field().String())
	if err != nil {
		return false
	}

	return userDate.After(time.Now().UTC())
}

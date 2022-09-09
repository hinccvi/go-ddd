package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		verr := err.(validator.ValidationErrors)
		field := verr[len(verr)-1].Field()
		tag := verr[len(verr)-1].Tag()
		msg := fmt.Sprintf("'%s' %s", field, tag)

		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}
	return nil
}

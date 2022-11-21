package middleware

import (
	"errors"
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
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			msg := fmt.Sprintf("'%s' %s", verr[len(verr)-1].Field(), verr[len(verr)-1].Tag())

			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		return err
	}

	return nil
}

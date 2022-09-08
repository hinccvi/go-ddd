package tools

import "github.com/labstack/echo/v4"

func Validator[I any](c echo.Context, i I) error {
	if err := c.Bind(&i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

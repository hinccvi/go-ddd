package tools

import "github.com/labstack/echo/v4"

func BindValidate[T any](c echo.Context, t T) error {
	if err := c.Bind(&t); err != nil {
		return err
	}

	return c.Validate(&t)
}

func JSONResp(c echo.Context, code int, msg message) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(msg),
	})
}

func JSONRespWithData[I any](c echo.Context, code int, msg message, i I) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(msg),
		Data:    i,
	})
}

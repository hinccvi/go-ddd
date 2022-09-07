package errors

import (
	"errors"
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
)

type (
	httpErrorHandler struct {
		statusCodes map[error]int
	}
)

func NewHttpErrorHandler(errorStatusCodeMaps map[error]int) *httpErrorHandler {
	return &httpErrorHandler{
		statusCodes: errorStatusCodeMaps,
	}
}

func unwrapRecursive(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)

		if internalErr == nil {
			break
		}

		originalErr = internalErr
	}

	return originalErr
}

func (eh *httpErrorHandler) GetStatusCode(err error) int {
	for key, value := range eh.statusCodes {
		if errors.Is(err, key) {
			return value
		}
	}

	return http.StatusInternalServerError
}

func (eh *httpErrorHandler) Handler(logger log.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    eh.GetStatusCode(err),
				Message: unwrapRecursive(err).Error(),
			}
		}

		code := he.Code
		message := he.Message
		if _, ok := he.Message.(string); ok {
			logger.With(c.Request().Context()).Error(err)
		}

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(he.Code)
			} else {
				err = tools.RespOkWithData(c, code, tools.MsgError, struct {
					Error string `json:"error"`
				}{message.(string)})
			}
			if err != nil {
				c.Echo().Logger.Error(err)
			}
		}
	}
}

package middleware

import (
	"errors"
	"net/http"

	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/labstack/echo/v4"
)

type HTTPErrorHandler struct {
	statusCodes map[error]int
}

func NewHTTPErrorHandler(errorStatusCodeMaps map[error]int) *HTTPErrorHandler {
	return &HTTPErrorHandler{
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

func (eh *HTTPErrorHandler) GetStatusCode(err error) int {
	for key, value := range eh.statusCodes {
		if errors.Is(err, key) {
			return value
		}
	}

	return http.StatusInternalServerError
}

//nolint:gocognit
func (eh *HTTPErrorHandler) Handler(logger log.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		var he *echo.HTTPError
		if errors.As(err, &he) {
			if he.Internal != nil {
				var ie *echo.HTTPError
				if errors.As(unwrapRecursive(he.Internal), &ie) {
					he = ie
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:     eh.GetStatusCode(err),
				Message:  "internal server error",
				Internal: err,
			}
		}

		l := logger.With(c.Request().Context(), "api", c.Request().RequestURI)
		code := he.Code
		message := ""

		if msg, ok := he.Message.(string); ok {
			message = msg
		} else {
			message = "invalid input"
		}

		if code == http.StatusInternalServerError {
			l.Error(message)
		}

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(he.Code)
			} else {
				err = tools.JSON(c, code, tools.Error, struct {
					Error string `json:"error"`
				}{message})
			}

			if err != nil {
				l.Error(he)
			}
		}
	}
}

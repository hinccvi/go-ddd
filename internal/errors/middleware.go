package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrBinding      = errors.New("invalid input")
	ErrInvalidToken = errors.New("invalid token")
	ErrUnauthorized = errors.New("unauthorized")
)

func Handler(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()
		if lastErr == nil {
			return
		}

		l := logger.With(c.Request.Context())

		l.Error(zap.Error(lastErr))

		c.JSON(http.StatusOK, buildErrorResponse(lastErr))
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err *gin.Error) ErrorResponse {
	switch e := err.Err.(type) {
	case validator.ValidationErrors:
		return InvalidInput(e[len(e)-1].Field() + " " + e[len(e)-1].Tag())
	case *json.SyntaxError:
		return InvalidInput("Invalid JSON format")
	case ErrorResponse:
		return e
	}

	if errors.Is(err, strconv.ErrSyntax) {
		return InvalidInput("Invalid syntax")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(err.Error())
	}

	return InternalServerError(err.Error())
}

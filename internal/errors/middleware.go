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
	ErrBinding  = errors.New("invalid input")
	ErrNotFound = errors.New("not found")
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

		r := buildErrorResponse(lastErr)
		c.JSON(http.StatusOK, r)
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err *gin.Error) ErrorResponse {
	if ve, ok := err.Err.(validator.ValidationErrors); ok {
		return InvalidInput(ve[len(ve)-1].Field() + " " + ve[len(ve)-1].Tag())
	}

	if _, ok := err.Err.(*json.SyntaxError); ok {
		return InvalidInput("Invalid JSON format")
	}

	// strconv failed
	if errors.Is(err, strconv.ErrSyntax) {
		return InvalidInput("Invalid syntax")
	}

	if errors.Is(err, ErrNotFound) {
		return NotFound(err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(err.Error())
	}

	return InternalServerError(err.Error())
}

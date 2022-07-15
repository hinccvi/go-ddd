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

		logger.With(c.Request.Context()).Error(zap.Error(lastErr))

		status := 0
		resp := buildErrorResponse(lastErr)

		switch resp.Status {
		case 400:
			fallthrough
		case 401:
			fallthrough
		case 403:
			fallthrough
		case 404:
			fallthrough
		case 500:
			status = resp.Status
		default:
			status = http.StatusBadRequest
		}

		c.JSON(status, resp)
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err *gin.Error) ErrorResponse {
	switch e := err.Err.(type) {
	case validator.ValidationErrors:
		return InvalidInput(404, e[len(e)-1].Field()+" "+e[len(e)-1].Tag())
	case *json.SyntaxError:
		return InvalidInput(404, "Invalid JSON format")
	case ErrorResponse:
		return e
	}

	if errors.Is(err, strconv.ErrSyntax) {
		return InvalidInput(404, "Invalid syntax")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(err.Error())
	}

	return InternalServerError()
}

package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type customValidator struct {
	Validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
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

func TestBindValidate_WhenOk(t *testing.T) {
	e := echo.New()
	e.Validator = &customValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"field":"bar"}`))
	req.Header.Set("Content-type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := struct {
		Field string `json:"field"`
	}{}

	assert.NoError(t, BindValidate(c, &data))
	assert.Equal(t, "bar", data.Field)
}

func TestBindValidate_WhenFailBind(t *testing.T) {
	e := echo.New()
	e.Validator = &customValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"field":"bar"`))
	req.Header.Set("Content-type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := struct {
		Text string `query:"field"`
	}{}

	assert.Error(t, BindValidate(c, &data))
}

func TestBindValidate_WhenFailValidate(t *testing.T) {
	e := echo.New()
	e.Validator = &customValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"field":""`))
	req.Header.Set("Content-type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := struct {
		Text string `query:"field" validator:"required"`
	}{}

	assert.Error(t, BindValidate(c, &data))
}

func TestResp_WithOk(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":200,"message":"success","data":null}`

	if assert.NoError(t, JSON(c, http.StatusOK, Success, nil)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestResp_WithCustomMessage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":200,"message":"test","data":null}`

	if assert.NoError(t, JSON(c, http.StatusOK, "test", nil)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestResp_WithBadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":400,"message":"error","data":null}`

	if assert.NoError(t, JSON(c, http.StatusBadRequest, Error, nil)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestResp_WithCustomError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":1000,"message":"error","data":null}`

	if assert.NoError(t, JSON(c, 1000, Error, nil)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestRespWithData_OkString(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":200,"message":"test","data":"string"}`

	if assert.NoError(t, JSON(c, http.StatusOK, "test", "string")) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestRespWithData_OkInt(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":200,"message":"test","data":1}`

	if assert.NoError(t, JSON(c, http.StatusOK, "test", 1)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestRespWithData_OkFloat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":200,"message":"test","data":1.1}`

	if assert.NoError(t, JSON(c, http.StatusOK, "test", 1.1)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestRespWithData_OkStruct(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "test",
		Age:  22,
	}

	b, err := json.Marshal(data)
	assert.Nil(t, err)

	okJSON := `{"code":200,"message":"test","data":` + string(b) + `}`

	if assert.NoError(t, JSON(c, http.StatusOK, "test", data)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

func TestRespWithData_CustomErrorStruct(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := struct {
		Error string `json:"error"`
	}{
		Error: "test",
	}

	b, err := json.Marshal(data)
	assert.Nil(t, err)

	okJSON := `{"code":1000,"message":"test","data":` + string(b) + `}`

	if assert.NoError(t, JSON(c, 1000, "test", data)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

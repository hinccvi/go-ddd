package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTestStatusCode(t *testing.T) {
	tests := []struct {
		code          int
		generatedCode int
	}{
		{code: http.StatusNotFound, generatedCode: http.StatusNotFound},
		{code: http.StatusInternalServerError, generatedCode: http.StatusInternalServerError},
		{code: http.StatusNetworkAuthenticationRequired, generatedCode: http.StatusNetworkAuthenticationRequired},
		{code: 9999, generatedCode: http.StatusBadRequest},
	}

	for _, test := range tests {
		assert.Equal(t, test.generatedCode, generateStatusCode(test.code))
	}
}

func TestResp_WithOk(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	okJSON := `{"code":200,"message":"success","data":null}`

	if assert.NoError(t, Resp(c, http.StatusOK, Success)) {
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

	if assert.NoError(t, Resp(c, http.StatusOK, "test")) {
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

	if assert.NoError(t, Resp(c, http.StatusBadRequest, Error)) {
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

	if assert.NoError(t, Resp(c, 1000, Error)) {
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

	if assert.NoError(t, RespWithData(c, http.StatusOK, "test", "string")) {
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

	if assert.NoError(t, RespWithData(c, http.StatusOK, "test", 1)) {
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

	if assert.NoError(t, RespWithData(c, http.StatusOK, "test", 1.1)) {
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

	if assert.NoError(t, RespWithData(c, http.StatusOK, "test", data)) {
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

	if assert.NoError(t, RespWithData(c, 1000, "test", data)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, okJSON, rec.Body.String())
	}
}

package http

import (
	"net/http"
	"testing"

	"github.com/hinccvi/go-ddd/internal/mocks"
	"github.com/hinccvi/go-ddd/internal/test"
	"github.com/hinccvi/go-ddd/pkg/log"
)

func TestAPI(t *testing.T) {
	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	router := mocks.Router(logger)

	RegisterHandlers(router.Group(""), "test")

	tests := []test.APITestCase{
		{
			Name:         "get",
			Method:       http.MethodGet,
			URL:          "/healthcheck",
			WantStatus:   http.StatusOK,
			WantResponse: `*OK test*`,
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}

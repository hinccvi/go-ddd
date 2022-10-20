package healthcheck

import (
	"net/http"
	"testing"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/test"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

func TestAPI(t *testing.T) {
	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	router := test.MockRouter(logger)

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

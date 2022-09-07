package test

// MockRoutingContext creates a gin context for testing handlers.
// func MockRoutingContext(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
// 	res := httptest.NewRecorder()
// 	if req.Header.Get("Content-Type") == "" {
// 		req.Header.Set("Content-Type", "application/json")
// 	}
// 	ctx, _ := echo.CreateTestContext(res)
// 	return ctx, res
// }

// MockRouter creates a gin router for testing APIs.
// func MockRouter(logger log.Logger) *gin.Engine {
// 	e := gin.Default()
// 	e.Use(
// 		accesslog.Handler(logger),
// 		errors.Handler(logger),
// 	)
// 	return e
// }

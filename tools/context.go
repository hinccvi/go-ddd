package tools

import (
	"context"
	"time"
)

func DefaultCancelContext() context.Context {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	return ctx
}

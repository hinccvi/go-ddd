package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	InvalidArgumentErr string
	InternalErr        string
	UnavailableErr     string
)

//nolint:gosec //false positive
const (
	InvalidCredential     InvalidArgumentErr = "Invalid credential"
	EmptyField            InvalidArgumentErr = "Empty field"
	InvalidToken          InvalidArgumentErr = "Invalid token"
	InvalidRefreshToken   InvalidArgumentErr = "Invalid refresh token"
	SystemError           InternalErr        = "System error"
	ConditionNotFulfilled UnavailableErr     = "Condition not fulfilled"
	MaxLoginAttempt       UnavailableErr     = "Max login attempt reached"
)

func (e InvalidArgumentErr) E() error {
	return status.Error(codes.InvalidArgument, string(e))
}

func (e InternalErr) E() error {
	return status.Error(codes.Internal, string(e))
}

func (e UnavailableErr) E() error {
	return status.Error(codes.Unavailable, string(e))
}

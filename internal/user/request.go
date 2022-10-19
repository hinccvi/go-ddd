package user

import "github.com/google/uuid"

type (
	getUserRequest struct {
		ID uuid.UUID `param:"id" validate:"required"`
	}

	queryUserRequest struct {
		Limit  int32  `query:"limit"`
		Offset *int32 `query:"offset"`
	}

	createUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	updateUserRequest struct {
		ID       uuid.UUID `json:"id" validate:"required"`
		Username string    `json:"username"`
		Password string    `json:"password"`
	}

	deleteUserRequest struct {
		ID uuid.UUID `param:"id" validate:"required"`
	}
)

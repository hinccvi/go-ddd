package user

import "github.com/google/uuid"

type (
	getOrDeleteUserRequest struct {
		ID uuid.UUID `query:"id" validate:"required"`
	}

	queryUserRequest struct {
		Limit  uint32  `query:"limit"`
		Offset *uint32 `query:"offset"`
	}

	createUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	updateUserRequest struct {
		ID       *uuid.UUID `json:"id" validate:"required"`
		Username string     `json:"username"`
		Password string     `json:"password"`
	}
)

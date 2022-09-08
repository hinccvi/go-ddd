package user

import "github.com/google/uuid"

type (
	getOrDeleteUserRequest struct {
		Id *uuid.UUID `form:"id" validate:"required"`
	}

	queryUserRequest struct {
		Limit  int  `form:"limit" validate:"required"`
		Offset *int `form:"offset" validate:"required"`
	}

	createUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	updateUserRequest struct {
		Id       string `json:"id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

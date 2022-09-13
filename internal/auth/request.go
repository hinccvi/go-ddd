package auth

type (
	loginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	refreshRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)

package auth

type (
	loginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	refreshResponse struct {
		RefreshToken string `json:"refresh_token"`
	}
)

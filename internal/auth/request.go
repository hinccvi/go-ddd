package auth

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

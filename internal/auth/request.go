package auth

type loginRequest struct {
	Name     string `binding:"required" json:"name"`
	Password string `binding:"required" json:"password"`
}

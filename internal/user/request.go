package user

type (
	getOrDeleteUserRequest struct {
		Id string `form:"id" validate:"required" `
	}

	queryUserRequest struct {
		Limit  int  `binding:"required" form:"limit"`
		Offset *int `binding:"required" form:"offset"`
	}

	createUserRequest struct {
		Name     string `binding:"required" json:"name"`
		Password string `binding:"required" json:"password"`
	}

	updateUserRequest struct {
		Id       string `binding:"required" json:"id"`
		Name     string `binding:"required" json:"name"`
		Password string `binding:"required" json:"password"`
	}
)

package user

type getOrDeleteUserRequest struct {
	Id string `binding:"required" form:"id"`
}

type queryUserRequest struct {
	Limit  int `binding:"required,gt=0" form:"limit"`
	Offset int `binding:"required,gt=0" form:"offset"`
}

type createUserRequest struct {
	Name     string `binding:"required" json:"name"`
	Age      int    `binding:"required" json:"age"`
	Position string `binding:"required" json:"position"`
}

type updateUserRequest struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Position string `json:"position"`
}

package user

type getOrDeleteUserRequest struct {
	Id string `binding:"required" form:"id"`
}

type queryUserRequest struct {
	Limit  int  `binding:"required,gt=0" form:"limit"`
	Offset *int `binding:"required,gte=0" form:"offset"`
}

type createUserRequest struct {
	Name     string `binding:"required" json:"name"`
	Password string `binding:"required" json:"password"`
}

type updateUserRequest struct {
	Id       string `binding:"required" json:"id"`
	Name     string `binding:"required" json:"name"`
	Password string `binding:"required" json:"password"`
}

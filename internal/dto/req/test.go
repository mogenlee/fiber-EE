package req

type TestListReq struct {
	Username string `query:"username"`
	Id       int32  `query:"id"`
}

type TestEditReq struct {
	Id       int32  `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type TestAddReq struct {
	Username string `json:"username" validate:"required"`
}

type TestLoginReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=18"`
}

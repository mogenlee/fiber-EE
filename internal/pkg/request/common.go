package request

// CommonIdReq 通用ID
type CommonIdReq struct {
	Id int32 `query:"id" form:"id" json:"id" validate:"required"`
}
type CommonKeyReq struct {
	Key string `query:"key" form:"key" json:"key" validate:"required"`
}

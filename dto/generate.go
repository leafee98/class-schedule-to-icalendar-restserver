package dto

type GenerateByPlanTokenReq struct {
	Token string `form:"token" binding:"required"`
}

type GenerateByPlanShareReq struct {
	ShareID int64 `form:"shareId" binding:"required"`
}

type GenerateRes struct {
	Content string `json:"content" binding:"required"`
}

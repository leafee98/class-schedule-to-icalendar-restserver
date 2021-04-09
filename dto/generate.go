package dto

type GenerateByPlanTokenReq struct {
	Token string `form:"token" binding:"required"`
}

type GenerateByPlanShareTokenReq struct {
	Token int64 `json:"token" form:"token" binding:"required"`
}

type GenerateRes struct {
	Content string `json:"content" binding:"required"`
}

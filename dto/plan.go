package dto

// Name     string `json:"name" binding:"required"`
// Remark   string `json:"remark" binding:"required"`
// ID       int64  `json:"id" binding:"required"`
// PlanID   int64  `json:"planId" binding:"required"`
// ConfigID int64  `json:"configId" binding:"required"`

type PlanCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

type PlanCreateRes struct {
	ID int64 `json:"id" binding:"required"`
}

type PlanAddConfigReq struct {
	PlanID   int64 `json:"planId" binding:"required"`
	ConfigID int64 `json:"configId" binding:"required"`
}

type PlanAddConfigRes struct {
	ID int64 `json:"id" binding:"required"`
}

type PlanRemoveConfigReq struct {
	PlanID   int64 `json:"planId" binding:"required"`
	ConfigID int64 `json:"configId" binding:"required"`
}

type PlanRemoveConfigRes string

type PlanCreateTokenReq struct {
	ID int64 `json:"id" binding:"required"`
}

type PlanCreateTokenRes struct {
	Token string `json:"token" binding:"required"`
}

type PlanRevokeTokenReq struct {
	Token string `json:"token" binding:"required"`
}

type PlanRevokeTokenRes string

// PlanGetTokenListReq is used to get all tokens of a plan
type PlanGetTokenListReq struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

// PlanGetTokenListRes is used to respond ConfigGetTokenListReq
type PlanGetTokenListRes struct {
	Tokens []PlanTokenDetail `json:"tokens" binding:"required"`
	Count  int64             `json:"count" binding:"required"`
}

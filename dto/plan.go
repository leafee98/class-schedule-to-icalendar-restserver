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

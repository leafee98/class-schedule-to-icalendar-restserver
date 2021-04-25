package dto

import "time"

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

type PlanAddConfigRes string

type PlanRemoveConfigReq struct {
	PlanID   int64 `json:"planId" binding:"required"`
	ConfigID int64 `json:"configId" binding:"required"`
}

type PlanRemoveConfigRes string

type PlanAddShareReq struct {
	PlanID        int64 `json:"planId" binding:"required"`
	ConfigShareID int64 `json:"configShareId" binding:"required"`
}

type PlanAddShareRes string

type PlanRemoveShareReq struct {
	PlanID        int64 `json:"planId" binding:"required"`
	ConfigShareID int64 `json:"configShareId" binding:"required"`
}

type PlanRemoveShareRes string

type PlanGetByIdReq struct {
	ID int64 `json:"id" binding:"required"`
}

type PlanGetByShareReq struct {
	ID int64 `json:"id" binding:"required"`
}

// response of PlanGetBy.. request
type PlanGetRes struct {
	Name       string         `json:"name" binding:"required"`
	Remark     string         `json:"remark" binding:"required"`
	ID         int64          `json:"id" binding:"required"`
	CreateTime time.Time      `json:"createTime" binding:"required"`
	ModifyTime time.Time      `json:"modifyTime" binding:"required"`
	Configs    []ConfigDetail `json:"configs" binding:"required"`

	// the share's detail is the same as configs, but its ID is shareID, not configID
	Shares []ConfigDetail `json:"shares" binding:"required"`
}

type PlanRemoveReq struct {
	ID int64 `json:"id" binding:"required"`
}

type PlanRemoveRes string

type PlanModifyReq struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark" binding:"required"`
}

type PlanModifyRes string

type PlanGetListReq struct {
	// available value: "createTime", "modifyTime", "name", "id"
	SortBy string `json:"sortBy" binding:"required"`

	// will make response's Count to zero when Offset is bigger than the number of config belongs to user
	Offset int64 `json:"offset"`

	// max 30
	Count int64 `json:"count" binding:"required"`
}

type PlanGetListRes struct {
	Count int64         `json:"count" binding:"required"`
	Plans []PlanSummary `json:"plans" binding:"required"`
}

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

type PlanShareCreateReq struct {
	ID     int64  `db:"c_id" json:"id" binding:"required"`
	Remark string `db:"c_remark" json:"remark" binding:"required"`
}

type PlanShareCreateRes struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

type PlanShareModifyReq struct {
	ID     int64  `db:"c_id" json:"id" binding:"required"`
	Remark string `db:"c_remark" json:"remark" binding:"required"`
}

type PlanShareModifyRes string

type PlanShareRevokeReq struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

type PlanShareRevokeRes string

type PlanShareGetListReq struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

type PlanShareGetListRes struct {
	Shares []PlanShareDetail `json:"shares" binding:"required"`
}

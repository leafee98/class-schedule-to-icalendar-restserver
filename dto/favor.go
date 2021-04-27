package dto

import "time"

type FavorConfigAddReq struct {
	ID int64 `json:"id" binding:"required"`
}
type FavorConfigAddRes string

type FavorConfigRemoveReq struct {
	ID int64 `json:"id" binding:"required"`
}
type FavorConfigRemoveRes string

type FavorConfigGetListReq struct {
	Offset int64 `json:"offset"`
	Count  int64 `json:"count" binding:"required"`
}

type FavorConfigGetListRes struct {
	Configs []FavorConfigSummary `json:"configs" binding:"required"`
}

type FavorPlanAddReq struct {
	ID int64 `json:"id" binding:"required"`
}
type FavorPlanAddRes string

type FavorPlanRemoveReq struct {
	ID int64 `json:"id" binding:"required"`
}
type FavorPlanRemoveRes string

type FavorPlanGetListReq struct {
	Offset int64 `json:"offset"`
	Count  int64 `json:"count" binding:"required"`
}

type FavorPlanGetListRes struct {
	Plans []FavorPlanSummary `json:"plans"`
}

/////////////////////////////////////
////////// Utility //////////////////
/////////////////////////////////////

type FavorConfigSummary struct {
	ShareID    int64     `json:"shareId" binding:"required"`
	Type       int8      `json:"type" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Format     int8      `json:"format" binding:"required"`
	Remark     string    `json:"remark" binding:"required"`
	FavorTime  time.Time `json:"favorTime" binding:"required"`
	CreateTime time.Time `json:"createTime" binding:"required"`
	ModifyTime time.Time `json:"modifyTime" binding:"required"`
}

type FavorPlanSummary struct {
	ShareID    int64     `json:"shareId" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Remark     string    `json:"remark" binding:"required"`
	FavorTime  time.Time `json:"favorTime" binding:"required"`
	CreateTime time.Time `json:"createTime" binding:"required"`
	ModifyTime time.Time `json:"modifyTime" binding:"required"`
}

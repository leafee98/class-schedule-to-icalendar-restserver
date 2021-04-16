package dto

import "time"

// ConfigSummary is used in http response
type ConfigSummary struct {
	ID         int64     `db:"c_id" json:"id" binding:"required"`
	Type       int8      `db:"c_type" json:"type" binding:"required"`
	Name       string    `db:"c_name" json:"name" binding:"required"`
	Format     int8      `db:"c_format" json:"format" binding:"required"`
	Remark     string    `db:"c_remark" json:"remark" binding:"required"`
	CreateTime time.Time `db:"c_create_time" json:"createTime" binding:"required"`
	ModifyTime time.Time `db:"c_modify_time" json:"modifyTime" binding:"required"`
}

// ConfigDetail is used in http response
type ConfigDetail struct {
	ID         int64     `db:"c_id" json:"id" binding:"required"`
	Type       int8      `db:"c_type" json:"type" binding:"required"`
	Name       string    `db:"c_name" json:"name" binding:"required"`
	Format     int8      `db:"c_format" json:"format" binding:"required"`
	Content    string    `db:"c_content" json:"content" binding:"required"`
	Remark     string    `db:"c_remark" json:"remark" binding:"required"`
	CreateTime time.Time `db:"c_create_time" json:"createTime" binding:"required"`
	ModifyTime time.Time `db:"c_modify_time" json:"modifyTime" binding:"required"`
}

type PlanTokenDetail struct {
	Token      string    `json:"tokens" binding:"required"`
	CreateTime time.Time `json:"createTime" binding:"required"`
}

type PlanSummary struct {
	ID         int64     `db:"c_id" json:"id" binding:"required"`
	Name       string    `db:"c_name" json:"name" binding:"required"`
	Remark     string    `db:"c_remark" json:"remark" binding:"required"`
	CreateTime time.Time `db:"c_create_time" json:"createTime" binding:"required"`
	ModifyTime time.Time `db:"c_modify_time" json:"modifyTime" binding:"required"`
}

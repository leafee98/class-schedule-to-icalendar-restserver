package dto

import "time"

// Full Database Properties
// ID         int64     `db:"c_id" json:"id" binding:"required"`
// Type       int8      `db:"c_type" json:"type" binding:"required"`
// Name       string    `db:"c_name" json:"name" binding:"required"`
// Content    string    `db:"c_content" json:"content" binding:"required"`
// Format     int8      `db:"c_format" json:"format" binding:"required"`
// OwnerID    int64     `db:"c_owner_id" json:"ownerId" binding:"required"`
// Remark     string    `db:"c_remark" json:"remark" binding:"required"`
// CreateTime time.Time `db:"c_create_time" json:"createTime" binding:"required"`
// ModifyTime time.Time `db:"c_modify_time" json:"modifyTime" binding:"required"`
// Deleted    bool      `db:"c_deleted" json:"deleted" binding:"required"`

const (
	LimitConfigFormatMin = 1
	LimitConfigFormatMax = 2
	LimitConfigTypeMin   = 1
	LimitConfigTypeMax   = 2
)

// ConfigCreateReq is used in the request to create a new global config
type ConfigCreateReq struct {
	Name    string `db:"c_name" json:"name" binding:"required"`
	Type    int8   `db:"c_type" json:"type" binding:"required"`
	Content string `db:"c_content" json:"content" binding:"required"`
	Format  int8   `db:"c_format" json:"format" binding:"required"`
	Remark  string `db:"c_remark" json:"remark" binding:"required"`
}

// ConfigCreateRes is used as response of ConfigCreateReq
type ConfigCreateRes struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

// ConfigGetByIDReq is used to get config's info with it's ID
type ConfigGetByIDReq struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

// ConfigGetByShareReq is used to get config's info with share link's ID
type ConfigGetByShareReq struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

// ConfigGetRes is used as response as ConfigGetReq
type ConfigGetRes struct {
	ID         int64     `db:"c_id" json:"id" binding:"required"`
	Type       int8      `db:"c_type" json:"type" binding:"required"`
	Name       string    `db:"c_name" json:"name" binding:"required"`
	Content    string    `db:"c_content" json:"content" binding:"required"`
	Format     int8      `db:"c_format" json:"format" binding:"required"`
	Remark     string    `db:"c_remark" json:"remark" binding:"required"`
	CreateTime time.Time `db:"c_create_time" json:"createTime" binding:"required"`
	ModifyTime time.Time `db:"c_modify_time" json:"modifyTime" binding:"required"`
}

// ConfigModifyReq is used to modify a config.
// Owner should not change the property Type.
type ConfigModifyReq struct {
	ID      int64  `db:"c_id" json:"id" binding:"required"`
	Name    string `db:"c_name" json:"name" binding:"required"`
	Content string `db:"c_content" json:"content" binding:"required"`
	Format  int8   `db:"c_format" json:"format" binding:"required"`
	Remark  string `db:"c_remark" json:"remark" binding:"required"`
}

// return status, succeed is "ok"
// Should be the same as Response.Status
type ConfigModifyRes string

// ConfigRemoveReq is used to remove config
type ConfigRemoveReq struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

// return status, succeed is "ok"
// Should be the same as Response.Status
type ConfigRemoveRes string

// ConfigGetListReq is used to get a list of user's config
type ConfigGetListReq struct {
	// available value: "createTime", "modifyTime", "name", "id"
	SortBy string `json:"sortBy" binding:"required"`

	// will make response's Count to zero when Offset is bigger than the number of config belongs to user
	Offset int64 `json:"offset"`

	// max 30
	Count int64 `json:"count" binding:"required"`
}

// ConfigGetListRes respond to ConfigGetListReq
type ConfigGetListRes struct {
	Count   int64           `json:"count" binding:"required"`
	Configs []ConfigSummary `json:"configs" binding:"required"`
}

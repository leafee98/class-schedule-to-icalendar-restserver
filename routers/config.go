package routers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/sirupsen/logrus"
)

func init() {
	RegisterRouter("/config-create", "post", configCreate)
	RegisterRouter("/config-get-by-id", "post", configGetByID)
	RegisterRouter("/config-modify", "post", configModify)
	RegisterRouter("/config-remove", "post", configRemove)
	RegisterRouter("/config-get-list", "post", configGetList)
}

// ConfigCreate will create a new Config
//
// Check if the user is authorized
// Check the type and format is in range of rule, err: "invalid type or format"
func configCreate(c *gin.Context) {
	var req dto.ConfigCreateReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// if the user is authorzied
	var ownerID int64
	if getUserIDOrAbort(c, &ownerID) != nil {
		return
	}

	// check config's type and format in range of rule
	if !checkConfigFormatRange(req.Format) || !checkConfigTypeRange(req.Type) {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("invalid type or format"))
		return
	}

	// insert the created config
	res, err := db.DB.Exec(
		"insert into t_config (c_type, c_name, c_content, c_format, c_owner_id, c_remark)"+
			" values (?, ?, ?, ?, ?, ?)",
		req.Type, req.Name, req.Content, req.Format, ownerID, req.Remark)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewResponseBad(err.Error()))
		return
	}

	configID, _ := res.LastInsertId()
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.ConfigCreateRes{ID: configID}))
}

// only owner could get config by id
//
// check login status, err msg: "unauthorized action is forbidden"
// check the deleted status, or no rows got, err msg: "config not exists"
// check the ownership, err msg: "you are not the owner of the config"
func configGetByID(c *gin.Context) {
	// bind request
	var req dto.ConfigGetByIDReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// check login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// get the config detail
	var res dto.ConfigGetRes
	var deleted bool
	var ownerID int64
	row := db.DB.QueryRow(
		"select c_id, c_type, c_name, c_content, c_format,"+
			"c_remark, c_create_time, c_modify_time,"+
			"c_deleted, c_owner_id from t_config where c_id = ?", req.ID)
	err := row.Scan(&res.ID, &res.Type, &res.Name, &res.Content, &res.Format,
		&res.Remark, &res.CreateTime, &res.ModifyTime,
		&deleted, &ownerID)

	if err != sql.ErrNoRows && err != nil {
		// unknown error
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
	} else if err == sql.ErrNoRows || deleted {
		// no such config or deleted
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("config not exists"))
	} else if err == nil {
		// check ownership
		if userID == ownerID {
			c.JSON(http.StatusOK, dto.NewResponseFine(res))
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("you are not the owner of the config"))
		}
	}
}

// other user could get config by share's id
//
// check the share link exists or not, err msg: "share link not found"
// check the share link expired status, err msg: "share link expired"
// check the config deleted status, err msg: "config not exists"
func configGetConfigByShare(c *gin.Context) {

}

// change config content
// property Type is not changable
//
// check login status, err msg: "unauthorized action is forbidden"
// check the deleted status, or no rows got, err msg: "config not exists"
// check the ownership, err msg: "you are not the owner of the config"
// update c_modify_time
func configModify(c *gin.Context) {
	// bind request
	var req dto.ConfigModifyReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// check login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// get the config detail
	var deleted bool
	var ownerID int64
	row := db.DB.QueryRow("select c_deleted, c_owner_id from t_config where c_id = ?", req.ID)
	err := row.Scan(&deleted, &ownerID)

	if err != sql.ErrNoRows && err != nil {
		// unknown error
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	} else if err == sql.ErrNoRows || deleted {
		// no such config or deleted
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("config not exists"))
		return
	} else if err == nil {
		// check ownership
		if userID != ownerID {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("you are not the owner of the config"))
			return
		}
	}

	// update the config
	res, err := db.DB.Exec("update t_config"+
		" set c_name=?, c_content=?, c_format=?, c_remark=?"+
		" where c_id = ?",
		req.Name, req.Content, req.Format, req.Remark, req.ID)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
	} else {
		affected, _ := res.RowsAffected()
		if affected > 0 {
			c.JSON(http.StatusOK, dto.NewResponseFine(dto.ConfigModifyRes("ok")))
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(dto.ConfigModifyRes("bad")))
		}
	}
}

// remove the config (set the deleted flag to true)
// check login
// check exists
// check ownership
func configRemove(c *gin.Context) {
	var req dto.ConfigRemoveReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check existence and ownership
	if configOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	const sqlCommand string = "update t_config set c_deleted = true where c_id = ?;"
	res, err := db.DB.Exec(sqlCommand, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.ConfigRemoveRes("ok")))
	} else {
		c.JSON(http.StatusBadGateway, dto.NewResponseFine(dto.ConfigRemoveRes("bad")))
	}
}

// get config's name, remark, create time and modify time owned by user.
// result will be from 'offset', 'count' no more than 30
func configGetList(c *gin.Context) {
	var req dto.ConfigGetListReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// count no more than 30
	if req.Count > 30 {
		req.Count = 30
	}

	switch req.SortBy {
	case "createTime":
		req.SortBy = "c_create_time"
	case "modifyTime":
		req.SortBy = "c_modify_time"
	case "name":
		req.SortBy = "c_name"
	case "id":
		req.SortBy = "c_id"
	default:
		req.SortBy = "c_id"
	}

	const sqlCommandPre = "select c_id, c_type, c_name, c_format, c_remark, c_create_time, c_modify_time from t_config" +
		" where c_owner_id = ? and c_deleted = false order by %s limit ?, ?;"
	var sqlCommand string = fmt.Sprintf(sqlCommandPre, req.SortBy)
	rows, err := db.DB.Query(sqlCommand, userID, req.Offset, req.Count)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error))
		return
	}

	var configSummarys []dto.ConfigSummary = make([]dto.ConfigSummary, 0)
	for rows.Next() {
		var configSummary dto.ConfigSummary
		rows.Scan(&configSummary.ID,
			&configSummary.Type,
			&configSummary.Name,
			&configSummary.Format,
			&configSummary.Remark,
			&configSummary.CreateTime,
			&configSummary.ModifyTime)

		configSummarys = append(configSummarys, configSummary)
	}
	rows.Close()

	c.JSON(http.StatusOK,
		dto.NewResponseFine(dto.ConfigGetListRes{Count: int64(len(configSummarys)), Configs: configSummarys}))
}

package routers

import (
	"database/sql"
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
	configID, err := res.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewResponseBad(err.Error()))
		return
	}
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

// get config's name, remark, create time and modify time owned by user.
// from 'start' to 'start + limit', 'limit' no more than 30
func configList(c *gin.Context) {

}

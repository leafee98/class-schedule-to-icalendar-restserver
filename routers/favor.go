package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/sirupsen/logrus"
)

// run automatically to register routers in this file
func init() {
	RegisterRouter("/favor-config-add", "post", favorConfigAdd)
	RegisterRouter("/favor-config-remove", "post", favorConfigRemove)
	RegisterRouter("/favor-config-get-list", "post", favorConfigGetList)

	RegisterRouter("/favor-plan-add", "post", favorPlanAdd)
	RegisterRouter("/favor-plan-remove", "post", favorPlanRemove)
	RegisterRouter("/favor-plan-get-list", "post", favorPlanGetList)
}

// check share existence
func favorConfigAdd(c *gin.Context) {
	var req dto.FavorConfigAddReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	if configShareExistOrAbort(c, req.ID) != nil {
		return
	}

	// check if the share is already in user's favor
	var cnt int64
	row := db.DB.QueryRow("select count(*) from t_user_favourite_config where c_config_share_id = ? and c_user_id = ?;",
		req.ID, userID)
	if err := row.Scan(&cnt); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}
	if cnt > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("this config share is already in your favor"))
		return
	}

	_, err := db.DB.Exec("insert into t_user_favourite_config (c_user_id, c_config_share_id) values (?, ?);",
		userID, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.FavorConfigAddRes("ok")))
}

func favorConfigRemove(c *gin.Context) {
	var req dto.FavorConfigRemoveReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	_, err := db.DB.Exec("delete from t_user_favourite_config where c_user_id = ? and c_config_share_id = ?;",
		userID, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.FavorConfigRemoveRes("ok")))
}

func favorConfigGetList(c *gin.Context) {
	var req dto.FavorConfigGetListReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	const sqlCommand string = `
		select
			tcs.c_id, tufc.c_create_time, tc.c_name, tc.c_remark, tc.c_type, tc.c_format, tc.c_create_time, tc.c_modify_time
		from t_config as tc 
			join t_config_share as tcs on tc.c_id = tcs.c_config_id
			join t_user_favourite_config as tufc on tcs.c_id = tufc.c_config_share_id
		where tc.c_deleted = false
			and tcs.c_deleted = false
			and tufc.c_user_id = ?
		limit ?, ?`

	rows, err := db.DB.Query(sqlCommand, userID, req.Offset, req.Count)
	defer rows.Close()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}

	var configs []dto.FavorConfigSummary = make([]dto.FavorConfigSummary, 0)
	for rows.Next() {
		var s dto.FavorConfigSummary
		err = rows.Scan(&s.ShareID, &s.FavorTime, &s.Name, &s.Remark, &s.Type, &s.Format, &s.CreateTime, &s.ModifyTime)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
			return
		}
		configs = append(configs, s)
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.FavorConfigGetListRes{Configs: configs}))
}

func favorPlanAdd(c *gin.Context) {
	var req dto.FavorPlanAddReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	if planShareExistOrAbort(c, req.ID) != nil {
		return
	}

	// check if the share is already in user's favor
	var cnt int64
	row := db.DB.QueryRow("select count(*) from t_user_favourite_plan where c_plan_share_id = ? and c_user_id = ?;",
		req.ID, userID)
	if err := row.Scan(&cnt); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}
	if cnt > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("this plan share is already in your favor"))
		return
	}

	_, err := db.DB.Exec("insert into t_user_favourite_plan (c_user_id, c_plan_share_id) values (?, ?);",
		userID, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.FavorPlanAddRes("ok")))
}

func favorPlanRemove(c *gin.Context) {
	var req dto.FavorPlanRemoveReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	_, err := db.DB.Exec("delete from t_user_favourite_plan where c_user_id = ? and c_plan_share_id = ?;",
		userID, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.FavorPlanRemoveRes("ok")))
}

func favorPlanGetList(c *gin.Context) {
	var req dto.FavorPlanGetListReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	const sqlCommand = `
		select tps.c_id, tufp.c_create_time, tp.c_name, tp.c_remark, tp.c_create_time, tp.c_modify_time
		from t_plan as tp 
			join t_plan_share as tps on tp.c_id = tps.c_plan_id
			join t_user_favourite_plan as tufp on tps.c_id = tufp.c_plan_share_id
		where tp.c_deleted = false
			and tps.c_deleted = false
			and tufp.c_user_id = ?
		limit ?, ?;`
	rows, err := db.DB.Query(sqlCommand, userID, req.Offset, req.Count)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
		return
	}

	var plans []dto.FavorPlanSummary = make([]dto.FavorPlanSummary, 0)
	for rows.Next() {
		var t dto.FavorPlanSummary
		err = rows.Scan(&t.ShareID, &t.FavorTime, &t.Name, &t.Remark, &t.CreateTime, &t.ModifyTime)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err))
			return
		}
		plans = append(plans, t)
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.FavorPlanGetListRes{Plans: plans}))
}

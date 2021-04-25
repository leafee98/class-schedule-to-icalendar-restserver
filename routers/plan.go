package routers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	RegisterRouter("plan-create", "post", planCreate)
	RegisterRouter("plan-add-config", "post", planAddConfig)
	RegisterRouter("plan-remove-config", "post", planRemoveConfig)
	RegisterRouter("plan-add-share", "post", planAddShare)
	RegisterRouter("plan-remove-share", "post", planRemoveShare)
	RegisterRouter("plan-get-by-id", "post", planGetById)
	RegisterRouter("plan-get-by-share", "post", planGetByShare)
	RegisterRouter("plan-remove", "post", planRemove)
	RegisterRouter("plan-modify", "post", planModify)
	RegisterRouter("plan-get-list", "post", planGetList)

	RegisterRouter("plan-create-token", "post", planCreateToken)
	RegisterRouter("plan-revoke-token", "post", planRevokeToken)
	RegisterRouter("plan-get-token-list", "post", planGetTokenList)

	RegisterRouter("/plan-share-create", "post", planShareCreate)
	RegisterRouter("/plan-share-modify", "post", planShareModify)
	RegisterRouter("/plan-share-revoke", "post", planShareRevoke)
	RegisterRouter("/plan-share-get-list", "post", planShareGetList)
}

// create a plan
// only need a Name and Remark, return ID
func planCreate(c *gin.Context) {
	// bind parameter
	var req dto.PlanCreateReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	res, err := db.DB.Exec("insert into t_plan (c_name, c_owner_id, c_remark) values (?, ?, ?);",
		req.Name, userID, req.Remark)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, err.Error())
		return
	}

	planID, _ := res.LastInsertId()
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanCreateRes{ID: planID}))
}

// check login status
// check if plan & config exist
// check if the relation already exist
// check ownership
// todo: transaction
func planAddConfig(c *gin.Context) {
	// bind request
	var req dto.PlanAddConfigReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// check login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check ownership
	if planOwnershipOrAbort(c, req.PlanID, userID) != nil {
		return
	}
	if configOwnershipOrAbort(c, req.ConfigID, userID) != nil {
		return
	}

	// check relation exist
	err3 := relationExist(req.PlanID, req.ConfigID)
	if err3 == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("this config already added to the plan"))
		return
	}

	// create relation
	_, err := db.DB.Exec("insert into t_plan_config_relation (c_plan_id, c_config_id) values (?, ?)",
		req.PlanID, req.ConfigID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanAddConfigRes("ok")))
}

func planRemoveConfig(c *gin.Context) {
	// bind request
	var req dto.PlanRemoveConfigReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// check login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check ownership
	if planOwnershipOrAbort(c, req.PlanID, userID) != nil {
		return
	}
	if configOwnershipOrAbort(c, req.ConfigID, userID) != nil {
		return
	}

	// check relation exist
	err3 := relationExist(req.PlanID, req.ConfigID)
	if err3 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("this config haven't been added to the plan"))
		return
	}

	// remove the relation
	const sqlCommand string = `delete from t_plan_config_relation where c_plan_id = ? and c_config_id = ?;`
	res, err := db.DB.Exec(sqlCommand, req.PlanID, req.ConfigID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		c.JSON(http.StatusOK, dto.NewResponseFine("ok"))
	} else {
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad("no relation deleted"))
	}
}

// check login status
// check if plan & config share exist
// check if the relation already exist
// check plan ownership
func planAddShare(c *gin.Context) {
	// bind request
	var req dto.PlanAddShareReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// check login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check plan ownership
	if planOwnershipOrAbort(c, req.PlanID, userID) != nil {
		return
	}
	// check config share existence
	if configShareExistOrAbort(c, req.ConfigShareID) != nil {
		return
	}

	// check relation exist
	err3 := relationShareExist(req.PlanID, req.ConfigShareID)
	if err3 == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("this config share already added to the plan"))
		return
	}

	// create relation
	_, err := db.DB.Exec("insert into t_plan_config_share_relation (c_plan_id, c_config_share_id) values (?, ?);",
		req.PlanID, req.ConfigShareID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanAddConfigRes("ok")))
}

func planRemoveShare(c *gin.Context) {
	// bind request
	var req dto.PlanRemoveShareReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// check login status
	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check plan ownership
	if planOwnershipOrAbort(c, req.PlanID, userID) != nil {
		return
	}

	// check relation of share exist
	err3 := relationShareExist(req.PlanID, req.ConfigShareID)
	if err3 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("this config haven't been added to the plan"))
		return
	}

	// remove the relation
	const sqlCommand string = `delete from t_plan_config_share_relation where c_plan_id = ? and c_config_share_id = ?;`
	res, err := db.DB.Exec(sqlCommand, req.PlanID, req.ConfigShareID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		c.JSON(http.StatusOK, dto.NewResponseFine("ok"))
	} else {
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad("no relation of share deleted"))
	}
}

// check login status
// check plan existence and ownership
func planGetById(c *gin.Context) {
	var req dto.PlanGetByIdReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	if planOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	var res dto.PlanGetRes
	if err := planGetRes(&res, req.ID); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
	} else {
		c.JSON(http.StatusOK, dto.NewResponseFine(res))
	}
}

func planGetByShare(c *gin.Context) {
	var req dto.PlanGetByShareReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	// get plan id
	var planID int64
	row := db.DB.QueryRow("select c_plan_id from t_plan_share where c_id = ?;", req.ID)
	if err := row.Scan(&planID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("plan share not exist"))
		return
	}

	var res dto.PlanGetRes
	if err := planGetRes(&res, planID); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
	} else {
		c.JSON(http.StatusOK, dto.NewResponseFine(res))
	}
}

func planRemove(c *gin.Context) {
	var req dto.PlanRemoveReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	if planOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	const sqlCommand string = "update t_plan set c_deleted = true where c_id = ?;"
	res, err := db.DB.Exec(sqlCommand, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	if affected, _ := res.RowsAffected(); affected > 0 {
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanRemoveRes("ok")))
	} else {
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad("deleted nothing"))
	}
}

func planModify(c *gin.Context) {
	var req dto.PlanModifyReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	if planOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	const sqlCommand = "update t_plan set c_name = ?, c_remark = ?  where c_id = ?;"
	_, err := db.DB.Exec(sqlCommand, req.Name, req.Remark, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanRemoveRes("ok")))
}

func planGetList(c *gin.Context) {
	var req dto.PlanGetListReq
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

	const sqlCommandPre = "select c_id, c_name, c_remark, c_create_time, c_modify_time" +
		" from t_plan where c_owner_id = ? and c_deleted = false order by %s limit ?, ?;"
	var sqlCommand string = fmt.Sprintf(sqlCommandPre, req.SortBy)
	rows, err := db.DB.Query(sqlCommand, userID, req.Offset, req.Count)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error))
		return
	}

	var planSummarys []dto.PlanSummary = make([]dto.PlanSummary, 0)
	for rows.Next() {
		var planSummary dto.PlanSummary
		rows.Scan(&planSummary.ID,
			&planSummary.Name,
			&planSummary.Remark,
			&planSummary.CreateTime,
			&planSummary.ModifyTime)

		planSummarys = append(planSummarys, planSummary)
	}
	rows.Close()

	c.JSON(http.StatusOK,
		dto.NewResponseFine(dto.PlanGetListRes{Count: int64(len(planSummarys)), Plans: planSummarys}))
}

// check delete status
// check ownership
func planCreateToken(c *gin.Context) {
	var req dto.PlanCreateTokenReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// existence and ownership
	if planOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	// cannot create token of plan more than 30
	var tokenCount int64
	row := db.DB.QueryRow("select count(*) from t_plan_token where c_plan_id = ?;", req.ID)
	err := row.Scan(&tokenCount)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
	}
	if tokenCount >= 30 {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(
			"number of tokens of the same plan cannot be more than 30, revoke some tokens before create more."))
		return
	}

	token := utils.GenerateToken()
	_, err = db.DB.Exec("insert into t_plan_token (c_plan_id, c_token) values (?, ?)", req.ID, token)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanCreateTokenRes{Token: token}))
}

// check login status
// check token existence
// check plan existence and ownership
func planRevokeToken(c *gin.Context) {
	var req dto.PlanRevokeTokenReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// get planID from token
	var planID int64
	row := db.DB.QueryRow("select c_plan_id from t_plan_token where c_token = ?;", req.Token)
	err := row.Scan(&planID)
	if err == sql.ErrNoRows {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("no such token"))
		return
	} else if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	if planOwnershipOrAbort(c, planID, userID) != nil {
		return
	}

	// delete this token
	res, err := db.DB.Exec("delete from t_plan_token where c_token = ?;", req.Token)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	if affected, _ := res.RowsAffected(); affected > 0 {
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanRevokeTokenRes("ok")))
	} else {
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(dto.PlanRevokeTokenRes("deleted nothing")))
	}
}

// check login status
// check plan existence and ownership
func planGetTokenList(c *gin.Context) {
	var req dto.PlanGetTokenListReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	if planOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	rows, err := db.DB.Query("select c_token, c_create_time from t_plan_token where c_plan_id = ?;", req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	var tokens []dto.PlanTokenDetail = make([]dto.PlanTokenDetail, 0)
	var token dto.PlanTokenDetail
	for rows.Next() {
		rows.Scan(&token.Token, &token.CreateTime)
		tokens = append(tokens, token)
	}
	rows.Close()
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanGetTokenListRes{Count: int64(len(tokens)), Tokens: tokens}))
}

func planShareCreate(c *gin.Context) {
	var req dto.PlanShareCreateReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check existence and ownership
	if planOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	const sqlCommand string = "insert into t_plan_share (c_plan_id, c_remark) values (?, ?);"
	res, err := db.DB.Exec(sqlCommand, req.ID, req.Remark)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	inserted, _ := res.LastInsertId()
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanShareCreateRes{ID: inserted}))
}

func planShareModify(c *gin.Context) {
	var req dto.PlanShareModifyReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check existence and ownership
	if planShareOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	const sqlCommand = "update t_plan_share set c_remark = ? where c_id = ?;"
	_, err := db.DB.Exec(sqlCommand, req.Remark, req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanShareModifyRes("ok")))
}

func planShareRevoke(c *gin.Context) {
	var req dto.PlanShareRevokeReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	var userID int64
	if getUserIDOrAbort(c, &userID) != nil {
		return
	}

	// check existence and ownership
	if planShareOwnershipOrAbort(c, req.ID, userID) != nil {
		return
	}

	_, err := db.DB.Exec("update t_plan_share set c_deleted = true where c_id = ?;", req.ID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanShareRevokeRes("ok")))
}

func planShareGetList(c *gin.Context) {
	var req dto.PlanShareGetListReq
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

	const sqlCommand = "select c_id, c_create_time, c_remark from t_plan_share " +
		"where c_deleted = false and c_plan_id = ?;"
	rows, err := db.DB.Query(sqlCommand, req.ID)
	defer rows.Close()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	var shareDetails []dto.PlanShareDetail = make([]dto.PlanShareDetail, 0)
	for rows.Next() {
		var shareDetail dto.PlanShareDetail
		if err := rows.Scan(&shareDetail.ID, &shareDetail.CreateTime, &shareDetail.Remark); err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
			return
		}
		shareDetails = append(shareDetails, shareDetail)
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanShareGetListRes{Shares: shareDetails}))
}

////////////////////////////////////////////////
/////////////////// Utilities //////////////////
////////////////////////////////////////////////

func planGetRes(plan *dto.PlanGetRes, planID int64) error {
	plan.Configs = make([]dto.ConfigDetail, 0)
	plan.Shares = make([]dto.ConfigDetail, 0)

	const sqlCommandGetPlan string = "select c_id, c_name, c_remark, c_create_time, c_modify_time " +
		"from t_plan where c_deleted = false and c_id = ?;"
	row := db.DB.QueryRow(sqlCommandGetPlan, planID)
	if err := row.Scan(&plan.ID, &plan.Name, &plan.Remark, &plan.CreateTime, &plan.ModifyTime); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("plan not exist")
		} else {
			return err
		}
	}

	const sqlCommandGetConfigs = "select " +
		"c_id, c_type, c_name, c_content, c_format, c_remark, c_create_time, c_modify_time " +
		"from t_config where c_id in (select c_config_id from t_plan_config_relation where c_plan_id = ?);"
	rows, err := db.DB.Query(sqlCommandGetConfigs, planID)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		var co dto.ConfigDetail
		err = rows.Scan(&co.ID, &co.Type, &co.Name, &co.Content, &co.Format, &co.Remark, &co.CreateTime, &co.ModifyTime)
		if err != nil {
			return err
		}
		plan.Configs = append(plan.Configs, co)
	}

	const sqlCommandGetShares = "" +
		"select s.c_id, c.c_type, c.c_name, c.c_content, c.c_format, c.c_remark, c.c_create_time, c.c_modify_time " +
		"from " +
		"	t_config as c " +
		"	join t_config_share as s " +
		"	on c.c_id = s.c_config_id " +
		"where " +
		"	c.c_deleted = false " +
		"	and s.c_deleted = false " +
		"	and	s.c_id in ( " +
		"		select c_config_share_id " +
		"		from t_plan_config_share_relation " +
		"		where c_plan_id = ? " +
		"	);"
	rowsShare, err := db.DB.Query(sqlCommandGetShares, planID)
	defer rowsShare.Close()
	if err != nil {
		return err
	}

	for rowsShare.Next() {
		var cs dto.ConfigDetail
		err = rowsShare.Scan(&cs.ID, &cs.Type, &cs.Name, &cs.Content, &cs.Format, &cs.Remark, &cs.CreateTime,
			&cs.ModifyTime)
		if err != nil {
			return err
		}
		plan.Shares = append(plan.Shares, cs)
	}

	return nil
}

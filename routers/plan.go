package routers

import (
	"database/sql"
	"errors"
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
	RegisterRouter("plan-get-by-id", "post", planGetById)
	RegisterRouter("plan-remove", "post", planRemove)
	RegisterRouter("plan-create-token", "post", planCreateToken)
	RegisterRouter("plan-revoke-token", "post", planRevokeToken)
	RegisterRouter("plan-get-token-list", "post", planGetTokenList)
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
	if planOwnerShipOrAbort(c, req.PlanID, userID) != nil {
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
	res, err := db.DB.Exec("insert into t_plan_config_relation (c_plan_id, c_config_id) values (?, ?)",
		req.PlanID, req.ConfigID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	relationID, _ := res.LastInsertId()
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanAddConfigRes{ID: relationID}))
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
	if planOwnerShipOrAbort(c, req.PlanID, userID) != nil {
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

	if planOwnerShipOrAbort(c, req.ID, userID) != nil {
		return
	}

	var res dto.PlanGetByIdRes
	res.Configs = make([]dto.ConfigDetail, 0)

	const sqlCommandGetPlan string = "select c_id, c_name, c_remark, c_create_time, c_modify_time " +
		"from t_plan where c_id = ?;"
	row := db.DB.QueryRow(sqlCommandGetPlan, req.ID)
	err := row.Scan(&res.ID, &res.Name, &res.Remark, &res.CreateTime, &res.ModifyTime)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	const sqlCommandGetConfigs = "select " +
		"c_id, c_type, c_name, c_content, c_format, c_remark, c_create_time, c_modify_time " +
		"from t_config where c_id in (select c_config_id from t_plan_config_relation where c_plan_id = ?);"
	rows, err := db.DB.Query(sqlCommandGetConfigs, req.ID)
	defer rows.Close()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	for rows.Next() {
		var co dto.ConfigDetail
		err = rows.Scan(&co.ID, &co.Type, &co.Name, &co.Content, &co.Format, &co.Remark, &co.CreateTime, &co.ModifyTime)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
			return
		}
		res.Configs = append(res.Configs, co)
	}

	c.JSON(http.StatusOK, dto.NewResponseFine(res))
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

	if planOwnerShipOrAbort(c, req.ID, userID) != nil {
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
	if planOwnerShipOrAbort(c, req.ID, userID) != nil {
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

	if planOwnerShipOrAbort(c, planID, userID) != nil {
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

	if planOwnerShipOrAbort(c, req.ID, userID) != nil {
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

///////////////////////////////
////// Utility Functions //////
///////////////////////////////

///////// Plan Part ///////////

// return nil if plan exist.
func planExist(planID int64) error {
	row := db.DB.QueryRow("select c_id from t_plan where c_id = ?", planID)
	return row.Scan(&planID)
}

func planExistOrAbort(c *gin.Context, planID int64) error {
	err := planExist(planID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("the plan not exist"))
		} else {
			logrus.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		}
	}
	return err
}

// return nil if the plan belongs to the user
func planOwnerShip(planID int64, userID int64) error {
	var userIDInDB int64
	row := db.DB.QueryRow("select c_owner_id from t_plan where c_id = ? and c_deleted = false", planID)
	err := row.Scan(&userIDInDB)
	if err == sql.ErrNoRows {
		return errors.New("the plan not exist")
	} else if err != nil {
		return err
	}
	if userIDInDB != userID {
		return errors.New("you are not the owner of the plan")
	}
	return nil
}

func planOwnerShipOrAbort(c *gin.Context, planID int64, userID int64) error {
	err := planOwnerShip(planID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	}
	return err
}

//////// Config Part //////////

// return nil if the config exist
func configExist(configID int64) error {
	row := db.DB.QueryRow("select c_id from t_plan where c_id = ? and c_deleted = false", configID)
	return row.Scan(&configID)
}

func configExistOrAbort(c *gin.Context, configID int64) error {
	err := configExist(configID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("the config not exist"))
		} else {
			logrus.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		}
	}
	return err
}

// return nil if the config belongs to the user
func configOwnership(configID int64, userID int64) error {
	var userIDInDB int64
	row := db.DB.QueryRow("select c_owner_id from t_config where c_id = ? and c_deleted = false", configID)
	err := row.Scan(&userIDInDB)
	if err == sql.ErrNoRows {
		return errors.New("the config not exist")
	} else if err != nil {
		return err
	}
	if userIDInDB != userID {
		return errors.New("you are not the owner of the plan")
	}
	return nil
}

func configOwnershipOrAbort(c *gin.Context, configID int64, userID int64) error {
	err := configOwnership(configID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	}
	return err
}

/////// Relation Part /////////

// return nil if exist
func relationExist(planID int64, configID int64) error {
	var cnt int64
	row := db.DB.QueryRow("select count(*) as cnt from t_plan_config_relation"+
		" where c_plan_id = ? and c_config_id = ?", planID, configID)
	err := row.Scan(&cnt)

	// sql.ErrNoRows will never occur
	if err != nil {
		return err
	} else {
		if cnt > 0 {
			return nil
		} else {
			return errors.New("no such relation")
		}
	}
}

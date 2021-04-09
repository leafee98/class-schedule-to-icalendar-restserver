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
	RegisterRouter("plan-create-token", "post", planCreateToken)
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

	token := utils.GenerateToken()
	res, err := db.DB.Exec("insert into t_plan_token (c_plan_id, c_token) values (?, ?)", req.ID, token)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}
	tokenID, _ := res.LastInsertId()
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.PlanCreateTokenRes{ID: tokenID, Token: token}))
}

// todo
func planRevokeToken(c *gin.Context) {

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
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
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

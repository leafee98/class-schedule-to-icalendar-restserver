package routers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/middlewares"
	"github.com/sirupsen/logrus"
)

func getUserIDOrAbort(c *gin.Context, userID *int64) error {
	idInterface, exist := c.Get(middlewares.Key.UserID)
	if exist == false {
		c.AbortWithStatusJSON(http.StatusForbidden,
			dto.NewResponseBad("unauthorized action is forbidden"))
		return errors.New("unanthorized")
	}
	*userID = idInterface.(int64)
	return nil
}

func checkConfigTypeRange(t int8) bool {
	return t <= dto.LimitConfigTypeMax && t >= dto.LimitConfigFormatMin
}

func checkConfigFormatRange(r int8) bool {
	return r <= dto.LimitConfigFormatMax && r >= dto.LimitConfigFormatMin
}

func bindOrAbort(c *gin.Context, req interface{}) error {
	err := c.ShouldBind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("invalid request parameters"))
	}
	return err
}

///////////////////////////////
/////// Database Utility //////
///////////////////////////////

///////// Plan Part ///////////

// return nil if plan exist.
func planExist(planID int64) error {
	row := db.DB.QueryRow("select c_id from t_plan where c_id = ? and c_deleted = false", planID)
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

func planOwnershipOrAbort(c *gin.Context, planID int64, userID int64) error {
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

///// Relation share Part /////

// return nil if exist
func relationShareExist(planID int64, configShareID int64) error {
	var cnt int64
	row := db.DB.QueryRow("select count(*) as cnt from t_plan_config_share_relation "+
		"where c_plan_id = ? and c_config_share_id = ?", planID, configShareID)
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

////// Config Share Part //////

func configShareExist(configShareID int64) error {
	var count int64
	row := db.DB.QueryRow("select count(*) from t_config_share where c_deleted = false and c_id = ?;", configShareID)
	if err := row.Scan(&count); err != nil {
		return err
	} else {
		if count > 0 {
			return nil
		} else {
			return errors.New("config share not exist or has been deleted")
		}
	}
}

func configShareExistOrAbort(c *gin.Context, configShareID int64) error {
	err := configShareExist(configShareID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	}
	return err
}

func configShareOwnership(configShareID int64, userID int64) error {
	var dbUserID int64
	const sqlCommand string = "select c_owner_id from t_config where c_deleted = false and c_id = " +
		"(select c_config_id from t_config_share where c_deleted = false and c_id = ?);"
	row := db.DB.QueryRow(sqlCommand, configShareID)
	err := row.Scan(&dbUserID)
	if err == sql.ErrNoRows {
		return errors.New("the config share doesn't exist")
	} else if err != nil {
		return err
	}
	if dbUserID != userID {
		return errors.New("you are not the owner of the config's share")
	}
	return nil
}

func configShareOwnershipOrAbort(c *gin.Context, configShareID int64, userID int64) error {
	err := configShareOwnership(configShareID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	}
	return err
}

/////// Plan Share Part ///////

func planShareExist(configShareID int64) error {
	var count int64
	row := db.DB.QueryRow("select count(*) from t_plan_share where c_deleted = false and c_id = ?;", configShareID)
	if err := row.Scan(&count); err != nil {
		return err
	} else {
		if count > 0 {
			return nil
		} else {
			return errors.New("plan share not exist or has been deleted")
		}
	}
}

func planShareExistOrAbort(c *gin.Context, configShareID int64) error {
	err := planShareExist(configShareID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	}
	return err
}

func planShareOwnership(configShareID int64, userID int64) error {
	var dbUserID int64
	const sqlCommand string = "select c_owner_id from t_plan where c_deleted = false and c_id = " +
		"(select c_plan_id from t_plan_share where c_deleted = false and c_id = ?);"
	row := db.DB.QueryRow(sqlCommand, configShareID)
	err := row.Scan(&dbUserID)
	if err == sql.ErrNoRows {
		return errors.New("the config share doesn't exist")
	} else if err != nil {
		return err
	}
	if dbUserID != userID {
		return errors.New("you are not the owner of the config's share")
	}
	return nil
}

func planShareOwnershipOrAbort(c *gin.Context, configShareID int64, userID int64) error {
	err := planShareOwnership(configShareID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	}
	return err
}

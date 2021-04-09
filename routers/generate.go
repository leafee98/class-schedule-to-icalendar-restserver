package routers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/rpc"
	"github.com/sirupsen/logrus"
)

func init() {
	RegisterRouter("/generate-by-plan", "get", generateByPlanToken)
}

// require the token in get request
// check the existence of token
// check the existence of plan
// validate config share
func generateByPlanToken(c *gin.Context) {
	var req dto.GenerateByPlanTokenReq
	if bindOrAbort(c, &req) != nil {
		return
	}

	const sqlGetConfig string = `select c_content, c_type, c_format from t_config 
		where c_id in (select c_config_id from t_plan_config_relation
			where c_plan_id in (
				select c_plan_id from t_plan_token where c_token = ?
			)
		)
		and c_deleted = false;`

	rows, err := db.DB.Query(sqlGetConfig, req.Token)
	defer rows.Close()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	const configJSONFormat string = `{ "global": %s, "lessons": [ %s ] }`
	var configGlobal string
	var configLessons []string = make([]string, 0)
	for rows.Next() {
		var content string
		var globalOrLesson, format int8
		err := rows.Scan(&content, &globalOrLesson, &format)
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if globalOrLesson == 1 {
			configGlobal = content
		} else if globalOrLesson == 2 {
			configLessons = append(configLessons, content)
		}
	}

	var configLessonsStr strings.Builder
	for i, s := range configLessons {
		if i != 0 {
			configLessonsStr.WriteString(",")
		}
		configLessonsStr.WriteString(s)
	}

	var configJSONRes string = fmt.Sprintf(configJSONFormat, configGlobal, configLessonsStr.String())

	generateRes, err := rpc.JSONGenerate(configJSONRes)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, generateRes)
}

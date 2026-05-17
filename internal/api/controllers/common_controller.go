package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/response"
)

type CommonController struct{}

var Common = &CommonController{}

func (c *CommonController) Health(ctx *gin.Context) {
	sqlDB, err := database.GetDB().DB()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusServiceUnavailable, response.ServiceCodeCommon, response.CaseCodeInternalError, nil, "database unavailable")
		return
	}
	if err := sqlDB.Ping(); err != nil {
		response.FailWithDetailed(ctx, http.StatusServiceUnavailable, response.ServiceCodeCommon, response.CaseCodeInternalError, nil, "database ping failed")
		return
	}
	response.Ok(ctx)
}

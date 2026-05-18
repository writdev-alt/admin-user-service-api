package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkgResponse "github.com/turahe/pkg/response"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/models/response"
	"github.com/writdev-alt/admin-user-service/internal/api/services"
)

type PermissionController struct {
	*BaseController
}

var Permission = &PermissionController{BaseController: NewBaseController()}

func (c *PermissionController) List(ctx *gin.Context) {
	var req request.PermissionListRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	rows, err := services.Permission.List(ctx.Request.Context(), req)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, p := range rows {
		data = append(data, response.ToPermissionResponse(p))
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, data, "Permissions retrieved successfully")
}

func (c *PermissionController) Detail(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	row, err := services.Permission.FindByUUID(ctx.Request.Context(), id)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	if row == nil {
		pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "Permission not found")
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeRetrieved, response.ToPermissionResponse(*row), "Permission retrieved successfully")
}

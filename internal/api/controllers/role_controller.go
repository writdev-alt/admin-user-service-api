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

type RoleController struct {
	*BaseController
}

var Role = &RoleController{BaseController: NewBaseController()}

func (c *RoleController) List(ctx *gin.Context) {
	var req request.RoleListRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	roles, total, err := services.Role.List(ctx.Request.Context(), req)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, r := range roles {
		data = append(data, response.ToRoleResponse(r))
	}
	pageNumber, pageSize := c.NormalizePagination(req.PageNumber, req.PageSize)
	pkgResponse.SimplePaginated(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, pkgResponse.SimplePaginationResponse{
		Data: data, PageNumber: pageNumber, PageSize: pageSize,
		HasNext: int64(pageNumber*pageSize) < total, HasPrev: pageNumber > 1,
	}, "Roles retrieved successfully")
}

func (c *RoleController) Detail(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	role, err := services.Role.FindByUUID(ctx.Request.Context(), id)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	if role == nil {
		pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "Role not found")
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeRetrieved, response.ToRoleResponse(*role), "Role retrieved successfully")
}

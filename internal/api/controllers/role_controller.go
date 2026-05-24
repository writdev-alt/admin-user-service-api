package controllers

import (
	"errors"
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
	roles, err := services.Role.List(ctx.Request.Context(), req)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, r := range roles {
		data = append(data, response.ToRoleResponse(r))
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, data, "Roles retrieved successfully")
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

func (c *RoleController) GetUsers(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	users, err := services.Role.GetUsers(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrRoleNotFound) {
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "Role not found")
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, u := range users {
		data = append(data, response.ToUserResponse(u))
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, data, "Role users retrieved successfully")
}

func (c *RoleController) GetPermissions(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	permissions, err := services.Role.GetPermissions(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrRoleNotFound) {
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "Role not found")
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, p := range permissions {
		data = append(data, response.ToPermissionResponse(p))
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, data, "Role permissions retrieved successfully")
}

func (c *RoleController) SetPermissions(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	var req request.SetRolePermissionsRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	role, err := services.Role.SetPermissions(ctx.Request.Context(), id, req.PermissionIDs)
	if err != nil {
		if errors.Is(err, services.ErrRoleNotFound) {
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "Role not found")
			return
		}
		if errors.Is(err, services.ErrPermissionNotFound) {
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "Permission not found")
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeUpdated, response.ToRoleResponse(*role), "Role permissions updated successfully")
}

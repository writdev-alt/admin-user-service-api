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

type UserController struct {
	*BaseController
}

var User = &UserController{BaseController: NewBaseController()}

func (c *UserController) List(ctx *gin.Context) {
	var req request.UserListRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	users, total, err := services.User.List(ctx.Request.Context(), req)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, u := range users {
		data = append(data, response.ToUserResponse(u))
	}
	pageNumber, pageSize := c.NormalizePagination(req.PageNumber, req.PageSize)
	pkgResponse.SimplePaginated(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, pkgResponse.SimplePaginationResponse{
		Data: data, PageNumber: pageNumber, PageSize: pageSize,
		HasNext: int64(pageNumber*pageSize) < total, HasPrev: pageNumber > 1,
	}, "Users retrieved successfully")
}

func (c *UserController) Detail(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	user, err := services.User.FindByUUID(ctx.Request.Context(), id)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	if user == nil {
		pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeRetrieved, response.ToUserResponse(*user), "User retrieved successfully")
}

func (c *UserController) Create(ctx *gin.Context) {
	var req request.UserCreateRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	user, err := services.User.Create(ctx.Request.Context(), req, 0)
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "username already taken" {
			pkgResponse.ValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.Created(ctx, pkgResponse.ServiceCodeCommon, response.ToUserResponse(*user), "User created successfully")
}

func (c *UserController) Update(ctx *gin.Context) {
	var req request.UserUpdateRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	user, err := services.User.Update(ctx.Request.Context(), id, req, 0)
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "username already taken" {
			pkgResponse.ValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	if user == nil {
		pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeRetrieved, response.ToUserResponse(*user), "User updated successfully")
}

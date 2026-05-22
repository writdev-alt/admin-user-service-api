package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/pkg/logger"
	pkgResponse "github.com/turahe/pkg/response"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/models/response"
	"github.com/writdev-alt/admin-user-service/internal/api/services"
)

type UserController struct {
	*BaseController
}

var User = &UserController{BaseController: NewBaseController()}

func (c *UserController) userResponse(ctx context.Context, u entities.User, roles map[uint64]*entities.Role) (response.UserResponse, error) {
	var role *entities.Role
	if roles != nil {
		role = roles[u.ID]
	} else {
		var err error
		role, err = services.User.RoleForUser(ctx, u.ID)
		if err != nil {
			return response.UserResponse{}, err
		}
	}
	return response.ToUserResponse(u, role), nil
}

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
	roles, err := services.User.RolesForUsers(ctx.Request.Context(), users)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	var data []interface{}
	for _, u := range users {
		resp, err := c.userResponse(ctx.Request.Context(), u, roles)
		if err != nil {
			pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
			return
		}
		data = append(data, resp)
	}
	pageNumber, pageSize := c.NormalizePagination(req.PageNumber, req.PageSize)
	pkgResponse.SimplePaginated(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeListRetrieved, pkgResponse.SimplePaginationResponse{
		Data: data, PageNumber: pageNumber, PageSize: pageSize,
		HasNext: int64(pageNumber*pageSize) < total, HasPrev: pageNumber > 1,
	}, "Users retrieved successfully")
}

func (c *UserController) Detail(ctx *gin.Context) {
	log := logger.WithContext(ctx.Request.Context())
	log.Infof("UserController.Detail: request_received")
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warnf("UserController.Detail: invalid id err=%v", err)
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	user, err := services.User.FindByUUID(ctx.Request.Context(), id)
	if err != nil {
		log.Warnf("UserController.Detail: error finding user err=%v", err)
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	if user == nil {
		log.Warnf("UserController.Detail: user not found")
		pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
		return
	}
	resp, err := c.userResponse(ctx.Request.Context(), *user, nil)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeRetrieved, resp, "User retrieved successfully")
}

func (c *UserController) Create(ctx *gin.Context) {
	log := logger.WithContext(ctx.Request.Context())
	log.Infof("UserController.Create: request_received")
	var req request.UserCreateRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		log.Warnf("UserController.Create: validation_failed err=%v", err)
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	adminID, ok := c.GetCurrentAdminID(ctx)
	if !ok {
		log.Warnf("UserController.Create: unauthorized")
		pkgResponse.UnauthorizedError(ctx, "Unauthorized")
		return
	}
	user, err := services.User.Create(ctx.Request.Context(), req, adminID)
	if err != nil {
		log.Warnf("UserController.Create: error creating user err=%v", err)
		if err.Error() == "email already registered" || err.Error() == "username already taken" {
			log.Warnf("UserController.Create: validation error err=%v", err)
			pkgResponse.ValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	log.Infof("UserController.Create: user created successfully uuid=%s", user.UUID)
	resp, err := c.userResponse(ctx.Request.Context(), *user, nil)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.Created(ctx, pkgResponse.ServiceCodeCommon, resp, "User created successfully")
}

func (c *UserController) Update(ctx *gin.Context) {
	log := logger.WithContext(ctx.Request.Context())
	log.Infof("UserController.Update: request_received")
	var req request.UserUpdateRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		log.Warnf("UserController.Update: validation_failed err=%v", err)
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warnf("UserController.Update: invalid id err=%v", err)
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	adminID, ok := c.GetCurrentAdminID(ctx)
	if !ok {
		log.Warnf("UserController.Update: unauthorized")
		pkgResponse.UnauthorizedError(ctx, "Unauthorized")
		return
	}
	user, err := services.User.Update(ctx.Request.Context(), id, req, adminID)
	if err != nil {
		log.Warnf("UserController.Update: error updating user err=%v", err)
		if err.Error() == "email already registered" || err.Error() == "username already taken" {
			pkgResponse.ValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	if user == nil {
		log.Warnf("UserController.Update: user not found")
		pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
		return
	}
	resp, err := c.userResponse(ctx.Request.Context(), *user, nil)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeRetrieved, resp, "User updated successfully")
}

func (c *UserController) Delete(ctx *gin.Context) {
	log := logger.WithContext(ctx.Request.Context())
	log.Infof("UserController.Delete: request_received")
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warnf("UserController.Delete: invalid id err=%v", err)
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	adminID, ok := c.GetCurrentAdminID(ctx)
	if !ok {
		log.Warnf("UserController.Delete: unauthorized")
		pkgResponse.UnauthorizedError(ctx, "Unauthorized")
		return
	}
	if err := services.User.Delete(ctx.Request.Context(), id, adminID); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			log.Warnf("UserController.Delete: user not found")
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeDeleted, nil, "User deleted successfully")
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	log := logger.WithContext(ctx.Request.Context())
	log.Infof("UserController.ChangePassword: request_received")
	var req request.ChangeUserPasswordRequest
	if err := c.ValidateReqParams(ctx, &req); err != nil {
		c.HandleValidationError(ctx, pkgResponse.ServiceCodeCommon, err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warnf("UserController.ChangePassword: invalid id err=%v", err)
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	adminID, ok := c.GetCurrentAdminID(ctx)
	if !ok {
		log.Warnf("UserController.ChangePassword: unauthorized")
		pkgResponse.UnauthorizedError(ctx, "Unauthorized")
		return
	}
	if err := services.User.ChangePassword(ctx.Request.Context(), id, req.NewPassword, adminID); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			log.Warnf("UserController.ChangePassword: user not found")
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
			return
		}
		log.Warnf("UserController.ChangePassword: error changing password err=%v", err)
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodePasswordChanged, gin.H{
		"message": "Password changed successfully",
	}, "Password changed successfully")
}

func (c *UserController) ToggleStatus(ctx *gin.Context) {
	log := logger.WithContext(ctx.Request.Context())
	log.Infof("UserController.ToggleStatus: request_received")
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Warnf("UserController.ToggleStatus: invalid id err=%v", err)
		pkgResponse.ValidationErrorSimple(ctx, pkgResponse.ServiceCodeCommon, "id", "The id must be a valid UUID.")
		return
	}
	adminID, ok := c.GetCurrentAdminID(ctx)
	if !ok {
		log.Warnf("UserController.ToggleStatus: unauthorized")
		pkgResponse.UnauthorizedError(ctx, "Unauthorized")
		return
	}
	user, err := services.User.ToggleStatus(ctx.Request.Context(), id, adminID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			log.Warnf("UserController.ToggleStatus: user not found")
			pkgResponse.NotFoundError(ctx, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeNotFound, "User not found")
			return
		}
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	resp, err := c.userResponse(ctx.Request.Context(), *user, nil)
	if err != nil {
		pkgResponse.FailWithDetailed(ctx, http.StatusInternalServerError, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeInternalError, nil, err.Error())
		return
	}
	pkgResponse.OkWithDetailed(ctx, http.StatusOK, pkgResponse.ServiceCodeCommon, pkgResponse.CaseCodeUpdated, resp, "User status toggled successfully")
}

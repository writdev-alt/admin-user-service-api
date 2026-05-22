package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/pkg/handler"
	"github.com/turahe/pkg/logger"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
)

type BaseController struct {
	handler.BaseHandler
	Repo *repositories.Repository
}

func NewBaseController() *BaseController {
	return &BaseController{Repo: repositories.Repo}
}

// GetCurrentUserUUID returns the authenticated principal UUID from the gin context.
func (c *BaseController) GetCurrentUserUUID(ctx *gin.Context) (uuid.UUID, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	switch v := userID.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, false
		}
		return parsed, true
	case uuid.UUID:
		return v, true
	default:
		return uuid.Nil, false
	}
}

// GetCurrentAdminID returns the authenticated admin's numeric ID (set by admin auth middleware).
func (c *BaseController) GetCurrentAdminID(ctx *gin.Context) (uint64, bool) {
	if v, ok := ctx.Get("admin_id"); ok {
		switch id := v.(type) {
		case uint64:
			return id, true
		case int:
			return uint64(id), true
		case int64:
			return uint64(id), true
		}
	}
	uid, ok := c.GetCurrentUserUUID(ctx)
	if !ok {
		return 0, false
	}
	admin, err := c.Repo.Admin.FindByUUID(ctx.Request.Context(), uid)
	if err != nil {
		logger.ErrorfContext(ctx.Request.Context(), "GetCurrentAdminID: FindByUUID(%s): %v", uid, err)
		return 0, false
	}
	if admin == nil {
		return 0, false
	}
	return admin.ID, true
}

package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/turahe/pkg/response"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type AdminFinder interface {
	FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Admin, error)
}

func RequireAdmin(repo AdminFinder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if v, ok := ctx.Get("is_impersonating"); ok {
			if impersonating, _ := v.(bool); impersonating {
				response.ForbiddenError(ctx, "Admin access required")
				ctx.Abort()
				return
			}
		}

		uid, ok := userUUIDFromContext(ctx)
		if !ok {
			response.UnauthorizedError(ctx, "Unauthorized")
			ctx.Abort()
			return
		}

		admin, err := repo.FindByUUID(ctx.Request.Context(), uid)
		if err != nil || admin == nil || !admin.Status {
			response.ForbiddenError(ctx, "Admin access required")
			ctx.Abort()
			return
		}

		ctx.Set("admin_id", admin.ID)
		ctx.Set("admin_uuid", admin.UUID.String())
		ctx.Next()
	}
}

func userUUIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	raw, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	switch v := raw.(type) {
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

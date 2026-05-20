package routes_v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/turahe/pkg/config"
	"github.com/turahe/pkg/jwt"
	pkgMiddlewares "github.com/turahe/pkg/middlewares"
	"github.com/writdev-alt/admin-user-service/internal/api/middleware"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
)

func Register(router *gin.Engine) {
	v1 := router.Group("")
	RegisterCommonRouter(v1.Group(""))

	verifier, err := jwt.NewVerifier(context.Background(), config.GetConfig())
	if err != nil {
		panic("JWT verifier required for protected routes: " + err.Error())
	}

	protected := v1.Group("")
	protected.Use(pkgMiddlewares.AuthMiddleware(verifier))
	protected.Use(middleware.RequireAdmin(repositories.Repo.Admin))
	protected.Use(middleware.RequireAdminSession())
	{
		RegisterUserRouter(protected.Group("users"))
		RegisterRoleRouter(protected.Group("roles"))
		RegisterPermissionRouter(protected.Group("permissions"))
	}
}

package routes_v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/turahe/pkg/config"
	"github.com/turahe/pkg/jwt"
	pkgMiddlewares "github.com/turahe/pkg/middlewares"
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
	{
		RegisterUserRouter(protected.Group("user"))
		RegisterRoleRouter(protected.Group("roles"))
		RegisterPermissionRouter(protected.Group("permissions"))
	}
}

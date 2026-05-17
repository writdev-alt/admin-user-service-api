package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	pkgMiddlewares "github.com/turahe/pkg/middlewares"
	routes_v1 "github.com/writdev-alt/admin-user-service/internal/api/routes/v1"
)

func Setup() *gin.Engine {
	app := gin.New()

	gin.DisableConsoleColor()
	gin.DefaultWriter = os.Stdout

	app.Use(pkgMiddlewares.TraceMiddleware())
	app.Use(pkgMiddlewares.LoggerMiddleware())
	app.Use(pkgMiddlewares.RecoveryHandler)
	app.Use(pkgMiddlewares.CORS())
	app.NoMethod(pkgMiddlewares.NoMethodHandler())
	app.NoRoute(pkgMiddlewares.NoRouteHandler())

	routes_v1.Register(app)
	return app
}

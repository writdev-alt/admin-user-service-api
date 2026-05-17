package routes_v1

import (
	controllers "github.com/writdev-alt/admin-user-service/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPermissionRouter(router *gin.RouterGroup) {
	router.GET("", controllers.Permission.List)
	router.GET("/:id", controllers.Permission.Detail)
}

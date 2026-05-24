package routes_v1

import (
	controllers "github.com/writdev-alt/admin-user-service/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRouter(router *gin.RouterGroup) {
	router.GET("", controllers.Role.List)
	router.GET("/:id/permissions", controllers.Role.GetPermissions)
	router.PUT("/:id/permissions", controllers.Role.SetPermissions)
	router.GET("/:id", controllers.Role.Detail)
}

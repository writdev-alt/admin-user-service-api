package routes_v1

import (
	controllers "github.com/writdev-alt/admin-user-service/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCommonRouter(router *gin.RouterGroup) {
	router.GET("/health", controllers.Common.Health)
}

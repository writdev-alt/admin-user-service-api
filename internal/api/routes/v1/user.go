package routes_v1

import (
	controllers "github.com/writdev-alt/admin-user-service/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	router.GET("", controllers.User.List)
	router.POST("", controllers.User.Create)
	router.GET("/:id", controllers.User.Detail)
	router.PUT("/:id", controllers.User.Update)
}

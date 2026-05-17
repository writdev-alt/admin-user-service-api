package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/turahe/pkg/config"
	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	"github.com/writdev-alt/admin-user-service/internal/api/routes"
)

func Run() {
	if err := config.Setup(""); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to setup config: %s\n", err)
		logger.Fatalf("failed to setup config, %s", err)
	}
	if err := database.Setup(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to setup database: %s\n", err)
		logger.Fatalf("failed to setup database, %s", err)
	}

	gin.SetMode(config.GetConfig().Server.Mode)

	port := config.GetConfig().Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	web := routes.Setup()

	fmt.Println("Admin User Service running on port " + port)
	fmt.Println("================================>")
	logger.Fatalf("%v", web.Run(":"+port))
}

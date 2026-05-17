package controllers

import (
	"github.com/turahe/pkg/handler"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
)

type BaseController struct {
	handler.BaseHandler
	Repo *repositories.Repository
}

func NewBaseController() *BaseController {
	return &BaseController{Repo: repositories.Repo}
}

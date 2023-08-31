package v0

import (
	"avito-trainee/internal/config"
	"avito-trainee/internal/service"
	"avito-trainee/internal/transport/rest/helpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	config   *config.Config
	helpers  *helpers.Manager
}

func NewHandler(services *service.Service, cfg *config.Config, helpers *helpers.Manager) *Handler {
	return &Handler{
		services: services,
		config:   cfg,
		helpers:  helpers,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v0 := api.Group("/v0")
	{
		h.initSegmentRoutes(v0)
		h.initUserRoutes(v0)
	}
}

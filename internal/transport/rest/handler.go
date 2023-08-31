package rest

import (
	"avito-trainee/internal/config"
	"avito-trainee/internal/service"
	"avito-trainee/internal/transport/rest/helpers"
	"avito-trainee/internal/transport/rest/middleware"
	"net/http"

	v0 "avito-trainee/internal/transport/rest/v0"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	cfg      *config.Config
	helpers  *helpers.Manager
}

func NewHandler(services *service.Service, cfg *config.Config, helpers *helpers.Manager) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
		helpers:  helpers,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	router := gin.New()

	if cfg.Env == config.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	prometheus := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)

	router.Use(
		middleware.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		middleware.Cors(),
		prometheus.Instrument(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	h.InitRoot(router)

	handlerV0 := v0.NewHandler(h.services, h.cfg, h.helpers)
	api := router.Group("/api")
	{
		handlerV0.Init(api)
	}
}

package helpers

import (
	"avito-trainee/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	VersionHeader       = "App-Version"
	SourceHeader        = "X-User-Source"
	ContentTypeHeader   = "Content-Type"
)

type Helpers interface {
	BindData(c *gin.Context, req interface{})
	ErrorsHandle(c *gin.Context, err error)
	LogError(err error)
	GetIdFromPath(c *gin.Context, key string) (uint, error)
}

type Manager struct {
	Logger *logger.LogManager
}

func NewManager(loggerManager *logger.LogManager) *Manager {
	return &Manager{
		Logger: loggerManager,
	}
}

package helpers

import (
	"avito-trainee/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func (m *Manager) BindData(c *gin.Context, req interface{}) error {
	err := c.Request.ParseForm()
	if err != nil {
		return err
	}

	if err := c.Bind(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument
			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}
			return models.ErrInvalidRequestParams
		}
		return err
	}

	return nil
}

func (m *Manager) GetIdFromPath(c *gin.Context, key string) (uint, error) {
	param := c.Param(key)
	intParam, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return uint(intParam), nil
}

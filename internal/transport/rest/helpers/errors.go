package helpers

import (
	"avito-trainee/internal/transport/rest/response"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "Minimum string length is " + fe.Param()
	case "max":
		return "Maximum string length is " + fe.Param()
	case "oneof":
		return "Field can be one of: " + fe.Param()
	}
	return "Unknown error"
}

func (m *Manager) LogError(err error) {
	m.Logger.Error(err)
}

func (m *Manager) ErrorsHandle(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.JsonResponse(c.Writer, response.Data{
			Code: http.StatusNotFound,
			Text: err.Error(),
		})
		return
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		response.JsonResponse(c.Writer, response.Data{
			Code: http.StatusBadRequest,
			Text: "Duplicate key value",
		})
		return
	}

	var validationErr validator.ValidationErrors

	if errors.As(err, &validationErr) {
		out := make([]ErrorMsg, len(validationErr))
		for i, vErr := range err.(validator.ValidationErrors) {
			out[i] = ErrorMsg{vErr.Field(), getErrorMsg(vErr)}
		}
		response.JsonResponse(c.Writer, response.Data{
			Code:         http.StatusBadRequest,
			ClientErrors: out,
		})
	} else {
		m.Logger.Error(err)
		response.JsonResponse(c.Writer, response.Data{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		})
	}
}

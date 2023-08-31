package v0

import (
	"avito-trainee/internal/dto/user"
	"avito-trainee/internal/transport/rest/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	userGroup := api.Group("/user")
	{
		userGroup.GET(":id", h.getUser)
		userGroup.POST("", h.createUser)
	}
}

func (h *Handler) createUser(c *gin.Context) {
	var createUserDto user.CreateUserDTO
	if err := h.helpers.BindData(c, &createUserDto); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	if err := createUserDto.Validate(); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	user, err := h.services.User.Create(createUserDto)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	response.JsonResponse(c.Writer, response.Data{
		Code: http.StatusOK,
		Data: user,
	})
}

func (h *Handler) getUser(c *gin.Context) {
	id, err := h.helpers.GetIdFromPath(c, "id")
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	user, err := h.services.User.Get(id)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	response.JsonResponse(c.Writer, response.Data{
		Code: http.StatusOK,
		Data: user,
	})
}

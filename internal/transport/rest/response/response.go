package response

import (
	"net/http"

	"github.com/alioygur/gores"
	"github.com/gin-gonic/gin"
)

type Data struct {
	Code         int         `json:"code,omitempty" default:"200"`
	Text         string      `json:"text,omitempty"`
	ClientErrors interface{} `json:"errors,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Meta         interface{} `json:"meta,omitempty"`
	Error        error       `json:"-"`
}

type response struct {
	Message string `json:"message"`
}

func NewResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{message})
}

func JsonResponse(w http.ResponseWriter, response Data) {
	err := gores.JSON(w, response.Code, response)
	if err != nil {
		return
	}
}

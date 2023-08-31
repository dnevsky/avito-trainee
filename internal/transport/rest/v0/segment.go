package v0

import (
	"avito-trainee/internal/dto/segment"
	"avito-trainee/internal/transport/rest/response"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initSegmentRoutes(api *gin.RouterGroup) {
	segmentGroup := api.Group("/segment")
	{
		segmentGroup.GET(":id", h.getSegment)
		segmentGroup.POST("", h.createSegment)
		segmentGroup.DELETE("", h.deleteSegment)
		segmentGroup.POST("/attach", h.attachSegments)
		segmentGroup.POST("/history", h.history)
	}
}

func (h *Handler) history(c *gin.Context) {
	var downloadHistory segment.DownloadHistoryDTO
	if err := h.helpers.BindData(c, &downloadHistory); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	file, err := h.services.Segment.DownloadHistory(downloadHistory)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	defer func() {
		fileName := file.Name()
		err := file.Close()
		if err != nil {
			return
		}
		os.Remove(fileName)
		if err != nil {
			return
		}
	}()

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}
}

func (h *Handler) attachSegments(c *gin.Context) {
	var attachSegments segment.AttachSegmentsDTO
	if err := h.helpers.BindData(c, &attachSegments); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	attached, deleted, err := h.services.Segment.AttachSegmentsToUser(attachSegments)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	response.JsonResponse(c.Writer, response.Data{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"deleted":  deleted,
			"attached": attached,
		},
	})
}

func (h *Handler) createSegment(c *gin.Context) {
	var createSegmentDto segment.CreateSegmentDTO
	if err := h.helpers.BindData(c, &createSegmentDto); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	if err := createSegmentDto.Validate(); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	segment, err := h.services.Segment.Create(createSegmentDto)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	response.JsonResponse(c.Writer, response.Data{
		Code: http.StatusOK,
		Data: segment,
	})
}

func (h *Handler) getSegment(c *gin.Context) {
	id := c.Param("id")

	segment, err := h.services.Segment.Get(id)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	response.JsonResponse(c.Writer, response.Data{
		Code: http.StatusOK,
		Data: segment,
	})
}

func (h *Handler) deleteSegment(c *gin.Context) {
	var deleteSegmentDto segment.DeleteSegmentDTO
	if err := h.helpers.BindData(c, &deleteSegmentDto); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	if err := deleteSegmentDto.Validate(); err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	err := h.services.Segment.Delete(deleteSegmentDto)
	if err != nil {
		h.helpers.ErrorsHandle(c, err)
		return
	}

	response.JsonResponse(c.Writer, response.Data{
		Code: http.StatusOK,
	})
}

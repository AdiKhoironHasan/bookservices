package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	*Handler
}

func NewBookHandler(h *Handler) *BookHandler {
	return &BookHandler{h}
}

func (r *BookHandler) List(c *gin.Context) {
	data, err := r.client.List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"data":    data.Books,
		"message": "Success get all books",
	})
}

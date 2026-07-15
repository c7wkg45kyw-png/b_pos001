package handler

import (
	"bpos001/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ok[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, model.APIResponse[T]{Success: true, Message: message, Data: data})
}
func created[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusCreated, model.APIResponse[T]{Success: true, Message: message, Data: data})
}
func noContent(c *gin.Context) {
	c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Message: "deleted"})
}
func fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, model.ErrorResponse{Success: false, Code: code, Message: message})
}
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if err == gorm.ErrRecordNotFound {
		fail(c, http.StatusNotFound, "NOT_FOUND", "record not found")
		return
	}
	fail(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
}

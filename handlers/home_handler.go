package handlers

import (
	"net/http"
	"todo-go-htmx/database"
	"todo-go-htmx/views"

	"github.com/gin-gonic/gin"
)

type HomeHtmlHandler struct {
	Repository *database.SQLiteRepository
}

func (h *HomeHtmlHandler) GetHome(ctx *gin.Context) {
	todoList, err := h.Repository.All()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "", views.Home(todoList))

}

package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"todo-go-htmx/database"
	"todo-go-htmx/models"
	"todo-go-htmx/views"

	"github.com/gin-gonic/gin"
)

type TodoHtmlHandler struct {
	Repository *database.SQLiteRepository
}

func (h *TodoHtmlHandler) GetAll(ctx *gin.Context) {
	todoList, err := h.Repository.All()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.HTML(http.StatusOK, "", views.TodoTable(todoList))
}

func (h *TodoHtmlHandler) Create(ctx *gin.Context) {
	description := ctx.PostForm("description")
	todo := models.Todo{
		Description: description,
		Completed:   false,
	}
	_, err := h.Repository.Create(todo)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	todoList, err := h.Repository.All()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.HTML(http.StatusOK, "", views.TodoTable(todoList))
}

func (h *TodoHtmlHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	todoToUpdate, err := h.Repository.GetById(int64(intId))
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	if completedValue := ctx.PostForm("completed"); completedValue != "" {
		todoToUpdate.Completed = !todoToUpdate.Completed
	}

	if description := ctx.PostForm("description"); description != "" {
		trimmedDescription := strings.TrimSpace(description)
		if trimmedDescription == "" {
			ctx.Header("HX-Reswap", "none")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Description cannot be empty"})
			return
		}
		todoMaxLength := 500
		if len(trimmedDescription) > todoMaxLength {
			ctx.Header("HX-Reswap", "none")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Description too long (max 500 characters)"})
			return
		}
		todoToUpdate.Description = trimmedDescription
	}

	updatedTodo, err := h.Repository.Update(int64(intId), *todoToUpdate)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "", views.TodoTableItem(*updatedTodo))
}

func (h *TodoHtmlHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//TODO: check if exists?
	err = h.Repository.Delete(int64(intId))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"todo-go-htmx/database"
	"todo-go-htmx/handlers"

	"github.com/gin-gonic/gin"
)

const fileName = "sqlite.db"

func main() {

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	repository := database.NewSQLiteRepository(db)

	if err := repository.Migrate(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migrated")

	router := gin.Default()
	router.HTMLRender = &TemplRender{}

	homeHtmlHandler := handlers.HomeHtmlHandler{Repository: repository}
	todoHtmlHandler := handlers.TodoHtmlHandler{Repository: repository}

	router.GET("/", homeHtmlHandler.GetHome)

	router.GET("/todo", todoHtmlHandler.GetAll)
	router.POST("/todo", todoHtmlHandler.Create)
	router.PATCH("/todo/:id", todoHtmlHandler.Update)
	router.DELETE("/todo/:id", todoHtmlHandler.Delete)

	router.Run(":8080")
}

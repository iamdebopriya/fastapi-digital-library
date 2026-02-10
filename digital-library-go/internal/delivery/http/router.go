package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *BookHandler, taskRunning *bool) {
	taskHandler := NewTaskHandler(taskRunning)
	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBookByID)
	r.POST("/books", h.CreateBook)
	r.PUT("/books/:id", h.UpdateBook)
	r.DELETE("/books/:id", h.DeleteBook)
	r.POST("/tasks/process", taskHandler.RunHeavyTask)
}

// @title Digital Library API
// @version 1.0
// @description Digital Library API migrated from FastAPI to Go using Gin.
// @host localhost:8080
// @BasePath /

package main

import (
	"log"
	"time"

	_ "digital-library-go/docs"
	"digital-library-go/internal/delivery/http"
	"digital-library-go/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*  GLOBAL TASK STATE  */

var taskRunning = false

/*  MIDDLEWARE: WAIT IF TASK RUNNING  */

func waitForTaskMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		for {
			http.TaskMu.Lock()
			running := taskRunning
			http.TaskMu.Unlock()

			if !running {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		c.Next()
	}
}

/*  X-PROCESS-TIME  */

type timingWriter struct {
	gin.ResponseWriter
	start time.Time
}

func (w timingWriter) WriteHeader(code int) {
	duration := time.Since(w.start)
	w.ResponseWriter.Header().Set("X-Process-Time", duration.String())
	w.ResponseWriter.WriteHeader(code)
}

func timingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Writer = timingWriter{ResponseWriter: c.Writer, start: start}
		c.Next()
	}
}

/*  CORS  */

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Expose-Headers", "X-Process-Time")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

/*  MAIN  */

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(waitForTaskMiddleware())
	r.Use(timingMiddleware())
	r.Use(corsMiddleware())

	uc := usecase.NewBookUsecase()
	bookHandler := http.NewBookHandler(uc)
	http.RegisterRoutes(r, bookHandler, &taskRunning)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server running on port 8080")
	r.Run(":8080")
}

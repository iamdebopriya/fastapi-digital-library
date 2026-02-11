// @title Digital Library API
// @version 1.0
// @description Digital Library API migrated from FastAPI to Go using Gin.
// @host localhost:8080
// @BasePath /

package main

import (
	"log"
	"time"

	_ "github.com/iamdebopriya/fastapi-digital-library/digital-library-go/docs"
	"github.com/iamdebopriya/fastapi-digital-library/digital-library-go/internal/delivery/http"
	"github.com/iamdebopriya/fastapi-digital-library/digital-library-go/internal/usecase"

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

/*  MIDDLEWARE: TIMING + USER-AGENT LOGGING  */
type timingWriter struct {
	gin.ResponseWriter
	start time.Time
}

func (w timingWriter) WriteHeader(code int) {
	duration := time.Since(w.start)
	w.ResponseWriter.Header().Set("X-Process-Time", duration.String())
	w.ResponseWriter.WriteHeader(code)
}

func timingAndUserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Log User-Agent
		userAgent := c.GetHeader("User-Agent")
		log.Printf("[LOG] Request received from: %s", userAgent)

		// Wrap writer for X-Process-Time
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

	// Middlewares
	r.Use(waitForTaskMiddleware())        // wait if task running
	r.Use(timingAndUserAgentMiddleware()) // X-Process-Time + log User-Agent
	r.Use(corsMiddleware())               // CORS

	// Book CRUD + Task Handlers
	uc := usecase.NewBookUsecase()
	bookHandler := http.NewBookHandler(uc)
	http.RegisterRoutes(r, bookHandler, &taskRunning)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server running on port 8080")
	r.Run(":8080")
}

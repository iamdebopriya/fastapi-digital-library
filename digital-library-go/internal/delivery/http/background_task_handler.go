package http

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

/* SHARED LOCK */
var TaskMu sync.Mutex

type TaskHandler struct {
	taskRunning *bool
}

func NewTaskHandler(flag *bool) *TaskHandler {
	return &TaskHandler{taskRunning: flag}
}

// RunHeavyTask godoc
// @Summary Run blocking background task
// @Description Runs a critical update task. Other requests wait until completion.
// @Tags Background Task
// @Produce json
// @Success 200 {object} map[string]string
// @Router /tasks/process [post]
func (h *TaskHandler) RunHeavyTask(c *gin.Context) {
	TaskMu.Lock()
	if *h.taskRunning {
		TaskMu.Unlock()
		c.JSON(http.StatusConflict, gin.H{"error": "task already running"})
		return
	}
	*h.taskRunning = true
	TaskMu.Unlock()

	log.Println("Task started")

	// Simulate heavy DB update
	time.Sleep(8 * time.Second)

	log.Println("Task finished")

	TaskMu.Lock()
	*h.taskRunning = false
	TaskMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Task completed successfully",
	})
}

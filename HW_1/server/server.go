package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "HW/docs"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Storage interface {
	Get(key string) (*string, error)
	Put(key string, value string) error
	Post(key string, value string) error
	Delete(key string) error
}

type Server struct {
	storage Storage
}

func newServer(storage Storage) *Server {
	return &Server{storage: storage}
}

type Task struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Compiler string `json:"compiler"`
	Status   string `json:"status"`
	Result   string `json:"result"`
}

// @Summary		Create a new task
// @Description	Creates a new task with given code and compiler.
// @Tags			task
// @Accept			json
// @Produce		json
// @Param			task	body		server.Task				true	"Task data"
// @Success		201		{object}	map[string]string		"Task created successfully"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/task [post]
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var taskData struct {
		Code     string `json:"code"`
		Compiler string `json:"compiler"`
	}
	if err := json.Unmarshal(body, &taskData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskID := uuid.New().String()
	task := Task{
		ID:       taskID,
		Code:     taskData.Code,
		Compiler: taskData.Compiler,
		Status:   "in_progress",
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.storage.Put(taskID, string(taskJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		time.Sleep(5 * time.Second)

		taskStr, err := s.storage.Get(taskID)
		if err != nil {
			log.Printf("Error getting task: %v", err)
			return
		}
		task := &Task{}
		if err := json.Unmarshal([]byte(*taskStr), task); err != nil {
			log.Printf("Error unmarshalling task: %v", err)
			return
		}
		task.Status = "ready"
		task.Result = "Task completed successfully."
		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Printf("Error marshalling task: %v", err)
			return
		}
		err = s.storage.Put(taskID, string(taskJSON))
		if err != nil {
			log.Printf("Error putting task: %v", err)
			return
		}
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})

	log.Printf("Sending response: %+v", taskID)
}

// @Summary		Get task status
// @Description	Gets the status of a task by its ID.
// @Tags			task
// @Produce		json
// @Param			task_id	query		string					true	"Task ID"
// @Success		200		{object}	map[string]string		"Task status"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		404		{object}	map[string]interface{}	"Task not found"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/status [get]
func (s *Server) getHandlerStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")

	if taskID == "" {
		http.Error(w, "task_id not provided", http.StatusBadRequest)
		return
	}

	taskStr, err := s.storage.Get(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	task := &Task{}
	if err := json.Unmarshal([]byte(*taskStr), task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": task.Status})
}

// @Summary		Get task result
// @Description	Gets the result of a task by its ID.
// @Tags			task
// @Produce		json
// @Param			task_id	query		string					true	"Task ID"
// @Success		200		{object}	map[string]string		"Task result"
// @Failure		400		{object}	map[string]interface{}	"Bad request"
// @Failure		404		{object}	map[string]interface{}	"Task not found or not ready"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/result [get]
func (s *Server) getHandlerResult(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")

	if taskID == "" {
		http.Error(w, "task_id not provided", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(taskID); err != nil {
		http.Error(w, "Invalid task_id parameter", http.StatusBadRequest)
		return
	}

	taskStr, err := s.storage.Get(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	task := &Task{}
	if err := json.Unmarshal([]byte(*taskStr), task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task.Status != "ready" {
		http.Error(w, "Task is not ready", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"result": task.Result})
}

func CreateAndRunServer(storage Storage, addr string) error {
	server := newServer(storage)

	r := chi.NewRouter()

	r.Post("/task", server.postHandler)
	r.Get("/status", server.getHandlerStatus)
	r.Get("/result", server.getHandlerResult)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/swagger/*", httpSwagger.WrapHandler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}

package dto

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
)

// TaskCreateRequest is the request body for the task create endpoint.
type TaskCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// TaskUpdateRequest is the request body for the task update endpoint.
type TaskUpdateRequest struct {
	ID     string       `json:"id"`
	Status types.Status `json:"status"`
}

// TaskResponse is the response body for the task.
type TaskResponse struct {
	ID          uint64       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      types.Status `json:"status"`
}

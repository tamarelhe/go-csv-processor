package api

import (
	"encoding/json"
	"net/http"

	"github.com/tamarelhe/go-csv-processor/services"
)

var progressService *services.ProgressService

// Initializes the API with the progress service
func InitProgressAPI(service *services.ProgressService) {
	progressService = service
}

// Returns the progress of the asynchronous upload
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	uploadID := r.URL.Query().Get("upload_id")
	if uploadID == "" {
		http.Error(w, "upload_id not provided", http.StatusBadRequest)
		return
	}

	// Get the upload progress through the service
	progress := progressService.GetProgress(uploadID)
	response := map[string]interface{}{
		"upload_id": uploadID,
		"progress":  progress,
	}

	// Returns progress in JSON format
	json.NewEncoder(w).Encode(response)
}

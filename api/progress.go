package api

import (
	"encoding/json"
	"net/http"

	"github.com/tamarelhe/go-csv-processor/services"
)

var progressService *services.ProgressService

// InitProgressAPI inicializa a API com o serviço de progresso
func InitProgressAPI(service *services.ProgressService) {
	progressService = service
}

// ProgressHandler retorna o progresso do upload assíncrono
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	uploadID := r.URL.Query().Get("upload_id")
	if uploadID == "" {
		http.Error(w, "upload_id não fornecido", http.StatusBadRequest)
		return
	}

	// Obtém o progresso do upload através do serviço
	progress := progressService.GetProgress(uploadID)
	response := map[string]interface{}{
		"upload_id": uploadID,
		"progress":  progress,
	}

	// Retorna o progresso em formato JSON
	json.NewEncoder(w).Encode(response)
}

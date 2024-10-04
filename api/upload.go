package api

import (
	"encoding/json"
	"net/http"

	"github.com/tamarelhe/go-csv-processor/services"
)

var uploadService *services.UploadService

func InitUploadAPI(service *services.UploadService) {
	uploadService = service
}

// UploadCSVHandler faz o upload de um CSV de forma genérica
func UploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("csv")
	if err != nil {
		http.Error(w, "Erro ao ler CSV", http.StatusBadRequest)
		return
	}
	defer file.Close()

	domain := r.FormValue("domain")
	uploadID := "some-upload-id" // Exemplo: você poderia usar um UUID aqui

	err = uploadService.Upload(domain, file, uploadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"upload_id": uploadID}
	json.NewEncoder(w).Encode(response)
}

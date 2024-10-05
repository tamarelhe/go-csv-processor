package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tamarelhe/go-csv-processor/domain"
	"github.com/tamarelhe/go-csv-processor/services"
)

// Structure to represent the download request with filters
type DownloadRequest struct {
	Domain  string          `json:"domain"`
	Filters []domain.Filter `json:"filters"`
}

var downloadService *services.DownloadService

// Initializes the API with the download service
func InitDownloadAPI(service *services.DownloadService) {
	downloadService = service
}

// Downloads the filtered CSV based on the filters sent in the body
func DownloadCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Only accepts POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decodes the JSON request body into the DownloadRequest structure
	var req DownloadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error processing the request body", http.StatusBadRequest)
		return
	}

	// Checks that the domain has been specified
	if req.Domain == "" {
		http.Error(w, "domain not provided", http.StatusBadRequest)
		return
	}

	// Calls the generic download service to generate the CSV
	csvData, err := downloadService.Download(req.Domain, req.Filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating CSV: %v", err), http.StatusInternalServerError)
		return
	}

	now := time.Now()
	formattedTime := now.Format("20060102_150405")
	content := fmt.Sprintf("attachment;filename=%s_%s.csv", req.Domain, formattedTime)

	// Configures HTTP header for CSV file download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", content)

	// Writes the CSV data in the body of the HTTP response
	w.Write(csvData)
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tamarelhe/go-csv-processor/domain"
	"github.com/tamarelhe/go-csv-processor/services"
)

// Estrutura para representar a requisição de download com filtros
type DownloadRequest struct {
	Domain  string          `json:"domain"`
	Filters []domain.Filter `json:"filters"`
}

var downloadService *services.DownloadService

// InitDownloadAPI inicializa a API com o serviço de download
func InitDownloadAPI(service *services.DownloadService) {
	downloadService = service
}

// DownloadCSVHandler faz o download do CSV filtrado com base nos filtros enviados no body
func DownloadCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Apenas aceita requisições POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o corpo da requisição JSON para a estrutura DownloadRequest
	var req DownloadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao processar o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Verifica se o domínio foi especificado
	if req.Domain == "" {
		http.Error(w, "Domínio não fornecido", http.StatusBadRequest)
		return
	}

	// Chama o serviço genérico de download para gerar o CSV
	csvData, err := downloadService.Download(req.Domain, req.Filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar CSV: %v", err), http.StatusInternalServerError)
		return
	}

	now := time.Now()
	formattedTime := now.Format("20060102_150405")
	content := fmt.Sprintf("attachment;filename=%s_%s.csv", req.Domain, formattedTime)

	// Configura o cabeçalho HTTP para download de arquivo CSV
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", content)

	// Escreve os dados do CSV no corpo da resposta HTTP
	w.Write(csvData)
}

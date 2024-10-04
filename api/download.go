package api

import (
	"fmt"
	"net/http"

	"github.com/tamarelhe/go-csv-processor/services"
)

var downloadService *services.DownloadService

// InitDownloadAPI inicializa a API com o serviço de download
func InitDownloadAPI(service *services.DownloadService) {
	downloadService = service
}

// DownloadCSVHandler faz o download do CSV filtrado
func DownloadCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Pega o domínio a partir dos parâmetros de consulta
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Domínio não fornecido", http.StatusBadRequest)
		return
	}

	// Constrói o mapa de filtros a partir dos parâmetros de consulta
	filters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if key != "domain" {
			filters[key] = values[0]
		}
	}

	// Chama o serviço genérico de download para gerar o CSV
	csvData, err := downloadService.Download(domain, filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Configura o cabeçalho HTTP para download de arquivo CSV
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=filtered_data.csv")

	// Escreve os dados do CSV no corpo da resposta HTTP
	w.Write(csvData)
}

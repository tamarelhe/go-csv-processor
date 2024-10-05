package main

import (
	"log"

	"net/http"

	"github.com/tamarelhe/go-csv-processor/test/files/purchase_order"

	"github.com/tamarelhe/go-csv-processor/api"
	"github.com/tamarelhe/go-csv-processor/domain"
	"github.com/tamarelhe/go-csv-processor/services"
)

func main() {

	// Inicializa os processadores específicos de domínio com as configurações desejadas
	po_processor := purchase_order.NewPOProcessor(
		true,
		';',
		[]domain.Column{
			{Label: "supplier", Type: domain.String, IsInputColumn: true, KeyColumn: true},
			{Label: "delivery_date", Type: domain.Date, IsInputColumn: true, KeyColumn: true},
			{Label: "item", Type: domain.String, IsInputColumn: true, KeyColumn: true},
			{Label: "location", Type: domain.Int, IsInputColumn: true, KeyColumn: true},
			{Label: "quantity", Type: domain.Float, IsInputColumn: true, KeyColumn: false},
		},
		true,
		false)

	processors := map[string]domain.CSVProcessor{
		"purchase_order": po_processor,
	}

	// Inicializa os serviços genéricos
	uploadService := services.NewUploadService(processors)
	downloadService := services.NewDownloadService(processors)
	progressService := services.NewProgressService()

	// Inicializa as APIs com os serviços
	api.InitUploadAPI(uploadService)
	api.InitDownloadAPI(downloadService)
	api.InitProgressAPI(progressService)

	// Configura as rotas HTTP
	http.HandleFunc("/upload", api.UploadCSVHandler)
	http.HandleFunc("/download", api.DownloadCSVHandler)
	http.HandleFunc("/progress", api.ProgressHandler)

	// Inicia o servidor HTTP
	log.Println("Servidor rodando em http://localhost:8091")
	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

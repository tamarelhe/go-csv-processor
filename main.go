package main

import (
	"log"

	"net/http"

	"github.com/tamarelhe/go-csv-processor/test/files/lead_time"
	"github.com/tamarelhe/go-csv-processor/test/files/purchase_order"

	"github.com/tamarelhe/go-csv-processor/api"
	"github.com/tamarelhe/go-csv-processor/domain"
	"github.com/tamarelhe/go-csv-processor/domain/models"
	"github.com/tamarelhe/go-csv-processor/services"
)

func main() {

	// Inicializa os processadores específicos de domínio com as configurações desejadas
	poProcessor := purchase_order.NewPOProcessor(
		true,
		';',
		[]models.Column{
			{Label: "supplier", Type: models.String, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "delivery_date", Type: models.Date, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "item", Type: models.String, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "location", Type: models.Int, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "quantity", Type: models.Float, IsMandatory: true, IsInputColumn: true, KeyColumn: false},
		},
		true,
		false)

	internalLeadTimeProcessor := lead_time.NewCfgInternalLeadTimeProcessor(
		true,
		';',
		[]models.Column{
			{Label: "operator", Type: models.String, IsMandatory: true, IsInputColumn: true, KeyColumn: false},
			{Label: "supplier", Type: models.String, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "location", Type: models.String, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "review_date", Type: models.Date, IsMandatory: true, IsInputColumn: true, KeyColumn: false},
			{Label: "review_day", Type: models.Int, IsMandatory: true, IsInputColumn: true, KeyColumn: true},
			{Label: "delivery_day", Type: models.Int, IsMandatory: true, IsInputColumn: true, KeyColumn: false},
			{Label: "min_intervel_weeks", Type: models.Int, IsMandatory: true, IsInputColumn: true, KeyColumn: false},
			{Label: "lead_time", Type: models.Int, IsMandatory: false, IsInputColumn: false, KeyColumn: false},
			{Label: "hash_control", Type: models.String, IsMandatory: false, IsInputColumn: true, KeyColumn: false},
		},
		true,
		true)

	processors := map[string]domain.CSVProcessor{
		"purchase_order":     poProcessor,
		"internal_lead_time": internalLeadTimeProcessor,
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

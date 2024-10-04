# go-csv-processor


## Using package im your own project

### Get package
```
go get github.com/tamarelhe/go-csv-processor
```

### Domain implementation
```
package delivery

import (
    "csv_processor/domain"
    "io"
    "fmt"
)

type DeliveryProcessor struct{}

func (p *DeliveryProcessor) ParseAndStage(file io.Reader, uploadID string) error {
    // Lógica de parsing específico para o domínio de entregas
    fmt.Println("Parse and stage delivery CSV")
    return nil
}

func (p *DeliveryProcessor) ApplyOperations() error {
    // Lógica para aplicar operações (Insert, Update, Delete)
    fmt.Println("Apply operations for delivery")
    return nil
}

func (p *DeliveryProcessor) GenerateCSV(filters map[string]string) ([]byte, error) {
    // Lógica para gerar o CSV filtrado
    return []byte("delivery;csv;data"), nil
}

```

### Using package
```
package main

import (
    "net/http"
    "csv_processor/api"
    "csv_processor/services"
    "delivery"
)

func main() {
    // Inicializa os processadores específicos de domínio
    processors := map[string]domain.CSVProcessor{
        "delivery": delivery.NewDeliveryProcessor(),
        // Adicionar outros domínios conforme necessário
    }

    // Inicializa os serviços genéricos de upload, download e progresso
    uploadService := services.NewUploadService(processors)
    downloadService := services.NewDownloadService(processors)
    progressService := services.NewProgressService()

    // Inicializa as APIs com os serviços
    api.InitUploadAPI(uploadService)
    api.InitDownloadAPI(downloadService)
    api.InitProgressAPI(progressService)

    // Configura as rotas de upload, download e progresso
    http.HandleFunc("/upload", api.UploadCSVHandler)
    http.HandleFunc("/download", api.DownloadCSVHandler)
    http.HandleFunc("/progress", api.ProgressHandler)

    // Inicia o servidor HTTP
    http.ListenAndServe(":8080", nil)
}
```

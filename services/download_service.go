package services

import (
	"fmt"

	"github.com/tamarelhe/go-csv-processor/domain"
)

// DownloadService gerencia o download de CSVs de forma genérica
type DownloadService struct {
	processors map[string]domain.CSVProcessor
}

// NewDownloadService inicializa o serviço com os processadores específicos
func NewDownloadService(processors map[string]domain.CSVProcessor) *DownloadService {
	return &DownloadService{
		processors: processors,
	}
}

// Download gera e retorna o CSV filtrado com base no domínio e filtros fornecidos
func (s *DownloadService) Download(domain string, filters map[string]string) ([]byte, error) {
	processor, exists := s.processors[domain]
	if !exists {
		return nil, fmt.Errorf("não há processador registrado para o domínio: %s", domain)
	}

	// Chama o método GenerateCSV do processador específico
	return processor.GenerateCSV(filters)
}

package services

import (
	"fmt"
	"io"
	"sync"

	"github.com/tamarelhe/go-csv-processor/domain"
)

// UploadService gerencia o upload de CSV de forma genérica
type UploadService struct {
	processors map[string]domain.CSVProcessor
	progress   map[string]int // Mapa para monitorar o progresso de uploads
	mu         sync.Mutex
}

// NewUploadService inicializa o serviço com processadores específicos
func NewUploadService(processors map[string]domain.CSVProcessor) *UploadService {
	return &UploadService{
		processors: processors,
		progress:   make(map[string]int),
	}
}

// Upload processa o CSV de acordo com o domínio e monitora o progresso
func (s *UploadService) Upload(domain string, file io.Reader, uploadID string) error {
	processor, exists := s.processors[domain]
	if !exists {
		return fmt.Errorf("não há processador registrado para o domínio: %s", domain)
	}

	// Simula progresso para fins de exemplo
	go func() {
		s.mu.Lock()
		s.progress[uploadID] = 0
		s.mu.Unlock()

		// Parsing e staging
		err := processor.ParseAndStage(file, uploadID)
		if err != nil {
			fmt.Println("Erro ao fazer o upload:", err)
		}

		// Simular progresso (100% completo)
		s.mu.Lock()
		s.progress[uploadID] = 100
		s.mu.Unlock()
	}()

	return nil
}

// GetProgress retorna o progresso de um upload
func (s *UploadService) GetProgress(uploadID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.progress[uploadID]
}

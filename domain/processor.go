package domain

import "io"

// CSVProcessor é a interface que o domínio específico deve implementar
type CSVProcessor interface {
	ParseAndStage(file io.Reader, uploadID string) error
	ApplyOperations() error
	GenerateCSV(filters map[string]string) ([]byte, error)
}

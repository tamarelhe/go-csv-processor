package domain

import "io"

type Processor struct {
	HasHeader bool
	Delimiter rune
}

// Filter representa um filtro aplicado no CSV
type Filter struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

// CSVProcessor é a interface que o domínio específico deve implementar
type CSVProcessor interface {
	ParseAndStage(file io.Reader, uploadID string) error
	ApplyOperations() error
	GenerateCSV(filters []Filter) ([]byte, error)
}

func NewProcessor(hasHeader bool, delimiter rune) *Processor {
	return &Processor{
		HasHeader: hasHeader,
		Delimiter: delimiter,
	}
}

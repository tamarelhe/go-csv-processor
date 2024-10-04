package domain

import "io"

// CSVFileDescriptor descreve a estrutura de um arquivo CSV
type CSVFileDescriptor struct {
	HasHeader bool     // Se o CSV tem cabeçalho
	Delimiter rune     // Delimitador do CSV
	Columns   []string // Nomes das colunas esperadas
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

func NewBaseCSVProcessor(hasHeader bool, delimiter rune) *CSVFileDescriptor {
	return &CSVFileDescriptor{
		HasHeader: hasHeader,
		Delimiter: delimiter,
	}
}

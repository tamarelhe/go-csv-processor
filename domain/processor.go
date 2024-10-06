package domain

import (
	"io"

	"github.com/tamarelhe/go-csv-processor/domain/models"
)

// Describes the structure of a CSV file
type CSVFileDescriptor struct {
	HasHeader          bool
	Delimiter          rune
	Columns            []models.Column
	ValidateUniqueness bool
	CUDControl         bool
}

// Represents a filter applied to the CSV
type Filter struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

// Interface that the specific domain must implement
type CSVProcessor interface {
	ParseAndStage(file io.Reader, uploadID string) error
	ApplyOperations() error
	GenerateCSV(filters []Filter) ([]byte, error)
	GetDescriptor() CSVFileDescriptor
}

func NewBaseCSVProcessor(hasHeader bool, delimiter rune, columns []models.Column, validateUniqueness bool, cudControl bool) *CSVFileDescriptor {
	return &CSVFileDescriptor{
		HasHeader:          hasHeader,
		Delimiter:          delimiter,
		Columns:            columns,
		ValidateUniqueness: validateUniqueness,
		CUDControl:         cudControl,
	}
}

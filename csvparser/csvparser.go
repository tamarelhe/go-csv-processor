package csvparser

import (
	"encoding/csv"
	"io"
)

// CSVParser estrutura para manter as configurações de parsing
type CSVParser struct {
	Delimiter rune
	HasHeader bool
}

// NewCSVParser inicializa um parser de CSV com delimitador e header
func NewCSVParser(delimiter rune, hasHeader bool) *CSVParser {
	return &CSVParser{
		Delimiter: delimiter,
		HasHeader: hasHeader,
	}
}

// Parse faz o parsing genérico do CSV
func (p *CSVParser) Parse(file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.Comma = p.Delimiter

	if p.HasHeader {
		if _, err := reader.Read(); err != nil {
			return nil, err
		}
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

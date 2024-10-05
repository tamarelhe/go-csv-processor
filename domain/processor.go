package domain

import (
	"io"
	"strconv"
	"time"
)

type ColumnType int

const (
	String ColumnType = iota
	Int
	Float
	Date
	DateTime
)

func (ct ColumnType) String() string {
	switch ct {
	case String:
		return "String"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Date:
		return "Date with 'DD/MM/YYYY' format'"
	case DateTime:
		return "DateTime with 'DD/MM/YYYY HH24:MI:SS' format"
	default:
		return "Unknown"
	}
}

const FormatDate = "02/01/2006"
const FormatDateTime = "02/01/2006 15:04:05"

type Column struct {
	Label         string
	Type          ColumnType
	IsInputColumn bool
	KeyColumn     bool
}

// Describes the structure of a CSV file
type CSVFileDescriptor struct {
	HasHeader          bool
	Delimiter          rune
	Columns            []Column
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

func NewBaseCSVProcessor(hasHeader bool, delimiter rune, columns []Column, validateUniqueness bool, cudControl bool) *CSVFileDescriptor {
	return &CSVFileDescriptor{
		HasHeader:          hasHeader,
		Delimiter:          delimiter,
		Columns:            columns,
		ValidateUniqueness: validateUniqueness,
		CUDControl:         cudControl,
	}
}

func IsValidDateFormat(dateStr string) bool {
	_, err := time.Parse(FormatDate, dateStr)

	return err == nil
}

func IsValidDateTimeFormat(dateTimeStr string) bool {
	_, err := time.Parse(FormatDateTime, dateTimeStr)

	return err == nil
}

func IsString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func IsStringInteger(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func IsStringFloat(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

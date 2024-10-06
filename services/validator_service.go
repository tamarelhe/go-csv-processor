package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tamarelhe/go-csv-processor/domain"
	"github.com/tamarelhe/go-csv-processor/domain/models"
)

const FormatDate = "02/01/2006"
const FormatDateTime = "02/01/2006 15:04:05"

func isValidDateFormat(dateStr string) bool {
	_, err := time.Parse(FormatDate, dateStr)

	return err == nil
}

func isValidDateTimeFormat(dateTimeStr string) bool {
	_, err := time.Parse(FormatDateTime, dateTimeStr)

	return err == nil
}

func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func isStringInteger(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func isStringFloat(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func isValidOperator(operator string) bool {
	switch models.Operator(operator) {
	case models.Add, models.Update, models.Delete:
		return true
	default:
		return false
	}
}

// Validate header length and columns position
func validateHeader(mapHeaders map[string]int, descriptor domain.CSVFileDescriptor) error {
	// Validates that the header matches the expected structure
	if len(mapHeaders) != len(descriptor.Columns) {
		return fmt.Errorf("number of columns in the header (%d) does not match the expected structure (%d)", len(mapHeaders), len(descriptor.Columns))
	}

	for _, col := range descriptor.Columns {
		_, exists := mapHeaders[col.Label]
		if !exists {
			return fmt.Errorf("column '%s' does not exists in the header", col.Label)
		}
	}

	// Validate operator column
	if descriptor.CUDControl {
		_, exists := mapHeaders["operator"]
		if !exists {
			return fmt.Errorf("column operator does not exists in the header")
		}

		_, exists = mapHeaders["hash_control"]
		if !exists {
			return fmt.Errorf("column hash_control does not exists in the header")
		}
	}

	return nil
}

// Validate line length and columns type
func validateRecord(headerMap map[string]int, recordNumber int, record []string, descriptor domain.CSVFileDescriptor) error {
	var columnIsValid = false

	// Validates that the record matches the expected structure
	if len(record) != len(descriptor.Columns) {
		return fmt.Errorf("record %d with incorrect number of columns: expected %d, but found %d", recordNumber, len(descriptor.Columns), len(record))
	}

	// Validate operator and hash control
	if descriptor.CUDControl {
		if !isValidOperator(record[headerMap["operator"]]) {
			return fmt.Errorf("invalid operator %s for record %d", record[0], recordNumber)
		}

		if record[headerMap["operator"]] == "A" && record[headerMap["hash_control"]] != "" {
			return fmt.Errorf("hash_control must be empty for add operator. record %d", recordNumber)
		}

		if (record[headerMap["operator"]] == "U" || record[headerMap["operator"]] == "D") && record[headerMap["hash_control"]] == "" {
			return fmt.Errorf("hash_control must be not empty for update or delete operator. record %d", recordNumber)
		}
	}

	// Validate columne types
	for _, colDef := range descriptor.Columns {
		if colDef.IsMandatory && record[headerMap[colDef.Label]] == "" {
			return fmt.Errorf("column %s of record %d cannot be empty", colDef.Label, recordNumber)
		}

		if record[headerMap[colDef.Label]] != "" {
			switch colDef.Type {
			case models.String:
				columnIsValid = isString(record[headerMap[colDef.Label]])
			case models.Int:
				columnIsValid = isStringInteger(record[headerMap[colDef.Label]])
			case models.Float:
				columnIsValid = isStringFloat(record[headerMap[colDef.Label]])
			case models.Date:
				columnIsValid = isValidDateFormat(record[headerMap[colDef.Label]])
			case models.DateTime:
				columnIsValid = isValidDateTimeFormat(record[headerMap[colDef.Label]])
			default:
				return fmt.Errorf("invalid type %s", colDef.Type)
			}

			if !columnIsValid {
				return fmt.Errorf("column %s of record %d has the value: %s. expected type %s", colDef.Label, recordNumber, record[headerMap[colDef.Label]], colDef.Type)
			}
		}
	}

	return nil
}

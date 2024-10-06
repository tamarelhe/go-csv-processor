package services

import (
	"bufio"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/tamarelhe/go-csv-processor/domain"
	"github.com/tamarelhe/go-csv-processor/domain/models"
)

type UploadService struct {
	processors map[string]domain.CSVProcessor
	state      map[string]models.State // Mapa para monitorar o progresso de uploads
	mu         sync.Mutex
}

func NewUploadService(processors map[string]domain.CSVProcessor) *UploadService {
	return &UploadService{
		processors: processors,
		state:      make(map[string]models.State),
	}
}

// Generates UploadID
func generateUploadID(domain string) string {
	timestamp := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(timestamp.UnixNano())), 0)

	ulidGenerated, err := ulid.New(ulid.Timestamp(timestamp), entropy)
	if err != nil {
		log.Fatal(err)
	}

	result := fmt.Sprintf("%s-%s", strings.ToUpper(domain), ulidGenerated.String())

	return result
}

// Generates UploadID
func generateRecordHash(record []string, descriptor domain.CSVFileDescriptor) ([]string, string, error) {
	var columnsKey []string

	for i, column := range descriptor.Columns {
		if column.KeyColumn {
			columnsKey = append(columnsKey, record[i])
		}
	}

	concatenatedString := strings.Join(columnsKey, " ")
	hash := sha256.Sum256([]byte(concatenatedString))
	hashString := hex.EncodeToString(hash[:])

	record = append(record, hashString)
	return record, hashString, nil
}

// Processes the CSV according to the domain and monitors progress
func (s *UploadService) Upload(domain string, file io.Reader) (string, error) {
	var headerMap map[string]int = make(map[string]int)
	uniqueHashes := make(map[string]int)

	uploadID := generateUploadID(domain)

	processor, exists := s.processors[domain]
	if !exists {
		return "", fmt.Errorf("there is no processor registered for the domain: %s", domain)
	}

	// Leitura inicial dos primeiros bytes do arquivo para verificar o BOM
	bufferedReader := bufio.NewReader(file)
	data, err := bufferedReader.Peek(3) // Pega os 3 primeiros bytes para verificar o BOM
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("erro ao ler arquivo: %v", err)
	}

	// Se o BOM for detectado, descarta esses 3 bytes
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		_, err := bufferedReader.Discard(3) // Descarta os 3 primeiros bytes do BOM
		if err != nil {
			return "", fmt.Errorf("erro ao descartar BOM: %v", err)
		}
	}

	// Reading the contents of the file to check for the presence of the BOM
	csvReader := csv.NewReader(bufferedReader)
	csvReader.Comma = processor.GetDescriptor().Delimiter
	csvReader.FieldsPerRecord = -1 // Allows records with a variable number of fields

	descriptor := processor.GetDescriptor()

	s.mu.Lock()
	s.state[uploadID] = models.Ready
	s.mu.Unlock()

	records, err := csvReader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("error reading the header: %v", err)
	}

	if len(records) == 0 {
		return "", fmt.Errorf("empty file")
	}

	if descriptor.HasHeader {
		headers := records[0]

		for i, header := range headers {
			headerMap[header] = i
		}

		if err = validateHeader(headerMap, descriptor); err != nil {
			return "", fmt.Errorf("invalid fields: %v", err)
		}
	}

	for i, record := range records[1:] {
		// Validates record
		if err = validateRecord(headerMap, i, record, descriptor); err != nil {
			return "", fmt.Errorf("invalid record [%d]: %v", i, err)
		}

		rec, hash, err := generateRecordHash(record, descriptor)
		if err != nil {
			return "", fmt.Errorf("invalid record [%d]: %v", i, err)
		}

		if descriptor.ValidateUniqueness {
			// Validates uniqueness of record
			if origRec, exists := uniqueHashes[hash]; exists {
				return "", fmt.Errorf("line %d is duplicated with line %d", i, origRec)
			}

			// Add hash to map
			uniqueHashes[hash] = i
		}

		fmt.Println(rec)
	}

	// Parsing & staging
	//err := processor.ParseAndStage(file, uploadID)
	//if err != nil {
	//	fmt.Println("Erro ao fazer o upload:", err)
	//}

	s.mu.Lock()
	s.state[uploadID] = models.Staged
	s.mu.Unlock()

	fmt.Printf("Processed %d records!", len(records[1:]))

	return uploadID, nil
}

// Returns the progress of an upload
func (s *UploadService) GetProgress(uploadID string) models.State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state[uploadID]
}

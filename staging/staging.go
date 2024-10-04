package staging

import "sync"

// StagingTable simula uma tabela de staging para armazenar dados temporariamente
type StagingTable struct {
	data map[string][]map[string]interface{} // Exemplo: pode ser uma interface mais sofisticada
	mu   sync.Mutex
}

// NewStagingTable cria uma nova tabela de staging
func NewStagingTable() *StagingTable {
	return &StagingTable{
		data: make(map[string][]map[string]interface{}),
	}
}

// Add adiciona um novo registro à tabela de staging
func (s *StagingTable) Add(uploadID string, record map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[uploadID] = append(s.data[uploadID], record)
}

// Get retorna os dados de staging para um upload específico
func (s *StagingTable) Get(uploadID string) []map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.data[uploadID]
}

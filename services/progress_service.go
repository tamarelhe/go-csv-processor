package services

// ProgressService monitora o progresso de uploads
type ProgressService struct {
	progress map[string]int
}

// NewProgressService inicializa o serviço de progresso
func NewProgressService() *ProgressService {
	return &ProgressService{
		progress: make(map[string]int),
	}
}

// GetProgress retorna o progresso de um upload
func (s *ProgressService) GetProgress(uploadID string) int {
	return s.progress[uploadID]
}

// UpdateProgress atualiza o progresso de um upload
func (s *ProgressService) UpdateProgress(uploadID string, percentage int) {
	s.progress[uploadID] = percentage
}

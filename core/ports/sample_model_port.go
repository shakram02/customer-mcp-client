package ports

import "mcp_client/core/domain"

// SampleModelPort defines the interface for sample model operations
type SampleModelPort interface {
	GetByID(id string) (*domain.SampleModel, error)
	Create(model *domain.SampleModel) error
	Update(model *domain.SampleModel) error
	Delete(id string) error
}

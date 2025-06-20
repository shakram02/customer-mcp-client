package sample_business_flow

import (
	"fmt"
	"mcp_client/core/domain"
	"mcp_client/core/ports"
)

// Input structs for the usecase
type GetSampleModelInput struct {
	ID string
}

type CreateSampleModelInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateSampleModelInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SampleBusinessFlowUsecase handles business logic for sample models
type SampleBusinessFlowUsecase struct {
	sampleModelRepo ports.SampleModelPort
}

// NewSampleBusinessFlowUsecase creates a new instance of the usecase
func NewSampleBusinessFlowUsecase(sampleModelRepo ports.SampleModelPort) *SampleBusinessFlowUsecase {
	return &SampleBusinessFlowUsecase{
		sampleModelRepo: sampleModelRepo,
	}
}

// GetSampleModel retrieves a sample model by ID
func (u *SampleBusinessFlowUsecase) GetSampleModel(input GetSampleModelInput) (*domain.SampleModel, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	model, err := u.sampleModelRepo.GetByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sample model: %w", err)
	}

	return model, nil
}

// CreateSampleModel creates a new sample model
func (u *SampleBusinessFlowUsecase) CreateSampleModel(input CreateSampleModelInput) error {
	if input.ID == "" || input.Name == "" {
		return fmt.Errorf("ID and Name are required")
	}

	model := domain.NewSampleModel(input.ID, input.Name, input.Description)

	if !model.IsValid() {
		return fmt.Errorf("invalid sample model data")
	}

	err := u.sampleModelRepo.Create(model)
	if err != nil {
		return fmt.Errorf("failed to create sample model: %w", err)
	}

	return nil
}

// UpdateSampleModel updates an existing sample model
func (u *SampleBusinessFlowUsecase) UpdateSampleModel(input UpdateSampleModelInput) error {
	if input.ID == "" || input.Name == "" {
		return fmt.Errorf("ID and Name are required")
	}

	model := domain.NewSampleModel(input.ID, input.Name, input.Description)

	if !model.IsValid() {
		return fmt.Errorf("invalid sample model data")
	}

	err := u.sampleModelRepo.Update(model)
	if err != nil {
		return fmt.Errorf("failed to update sample model: %w", err)
	}

	return nil
}

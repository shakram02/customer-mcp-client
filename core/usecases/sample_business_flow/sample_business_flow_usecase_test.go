package sample_business_flow

import (
	"errors"
	"mcp_client/core/domain"
	"testing"
)

// Mock implementation of SampleModelPort for testing
type mockSampleModelRepo struct {
	models      map[string]*domain.SampleModel
	shouldError bool
}

func (m *mockSampleModelRepo) GetByID(id string) (*domain.SampleModel, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}

	model, exists := m.models[id]
	if !exists {
		return nil, nil
	}
	return model, nil
}

func (m *mockSampleModelRepo) Create(model *domain.SampleModel) error {
	if m.shouldError {
		return errors.New("mock error")
	}

	if m.models == nil {
		m.models = make(map[string]*domain.SampleModel)
	}
	m.models[model.ID] = model
	return nil
}

func (m *mockSampleModelRepo) Update(model *domain.SampleModel) error {
	if m.shouldError {
		return errors.New("mock error")
	}

	if m.models == nil {
		m.models = make(map[string]*domain.SampleModel)
	}
	m.models[model.ID] = model
	return nil
}

func (m *mockSampleModelRepo) Delete(id string) error {
	if m.shouldError {
		return errors.New("mock error")
	}

	delete(m.models, id)
	return nil
}

func TestSampleBusinessFlowUsecase_GetSampleModel(t *testing.T) {
	tests := []struct {
		name        string
		input       GetSampleModelInput
		setupMock   func(*mockSampleModelRepo)
		expectError bool
		expectNil   bool
	}{
		{
			name:  "successful retrieval",
			input: GetSampleModelInput{ID: "test-id"},
			setupMock: func(m *mockSampleModelRepo) {
				m.models = map[string]*domain.SampleModel{
					"test-id": {ID: "test-id", Name: "Test", Description: "Test model"},
				}
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "empty ID",
			input:       GetSampleModelInput{ID: ""},
			setupMock:   func(m *mockSampleModelRepo) {},
			expectError: true,
			expectNil:   true,
		},
		{
			name:  "model not found",
			input: GetSampleModelInput{ID: "nonexistent"},
			setupMock: func(m *mockSampleModelRepo) {
				m.models = make(map[string]*domain.SampleModel)
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name:  "repository error",
			input: GetSampleModelInput{ID: "test-id"},
			setupMock: func(m *mockSampleModelRepo) {
				m.shouldError = true
			},
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSampleModelRepo{}
			tt.setupMock(mockRepo)

			usecase := NewSampleBusinessFlowUsecase(mockRepo)
			result, err := usecase.GetSampleModel(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if tt.expectNil && result != nil {
				t.Errorf("expected nil result but got: %v", result)
			}
			if !tt.expectNil && result == nil {
				t.Errorf("expected non-nil result but got nil")
			}
		})
	}
}

func TestSampleBusinessFlowUsecase_CreateSampleModel(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateSampleModelInput
		setupMock   func(*mockSampleModelRepo)
		expectError bool
	}{
		{
			name: "successful creation",
			input: CreateSampleModelInput{
				ID:          "test-id",
				Name:        "Test",
				Description: "Test model",
			},
			setupMock:   func(m *mockSampleModelRepo) {},
			expectError: false,
		},
		{
			name: "empty ID",
			input: CreateSampleModelInput{
				ID:   "",
				Name: "Test",
			},
			setupMock:   func(m *mockSampleModelRepo) {},
			expectError: true,
		},
		{
			name: "empty name",
			input: CreateSampleModelInput{
				ID:   "test-id",
				Name: "",
			},
			setupMock:   func(m *mockSampleModelRepo) {},
			expectError: true,
		},
		{
			name: "repository error",
			input: CreateSampleModelInput{
				ID:   "test-id",
				Name: "Test",
			},
			setupMock: func(m *mockSampleModelRepo) {
				m.shouldError = true
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSampleModelRepo{}
			tt.setupMock(mockRepo)

			usecase := NewSampleBusinessFlowUsecase(mockRepo)
			err := usecase.CreateSampleModel(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

package domain

// SampleModel represents a sample entity in our domain
type SampleModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewSampleModel creates a new SampleModel instance
func NewSampleModel(id, name, description string) *SampleModel {
	return &SampleModel{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

// IsValid validates the SampleModel
func (s *SampleModel) IsValid() bool {
	return s.ID != "" && s.Name != ""
}

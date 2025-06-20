package api

import (
	"encoding/json"
	"mcp_client/core/usecases/sample_business_flow"
	"net/http"
)

type SampleHandler struct {
	sampleUsecase *sample_business_flow.SampleBusinessFlowUsecase
}

func NewSampleHandler(sampleUsecase *sample_business_flow.SampleBusinessFlowUsecase) *SampleHandler {
	return &SampleHandler{
		sampleUsecase: sampleUsecase,
	}
}

func (h *SampleHandler) GetSampleModel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	input := sample_business_flow.GetSampleModelInput{
		ID: id,
	}

	model, err := h.sampleUsecase.GetSampleModel(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if model == nil {
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func (h *SampleHandler) CreateSampleModel(w http.ResponseWriter, r *http.Request) {
	var input sample_business_flow.CreateSampleModelInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.sampleUsecase.CreateSampleModel(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Sample model created successfully"))
}

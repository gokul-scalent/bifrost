package handlers

import (
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
)

func TestShouldReturnEmptyListModelsResponseForUnsupportedOperation(t *testing.T) {
	errCode := "unsupported_operation"
	err := &schemas.BifrostError{
		Error: &schemas.ErrorField{Code: &errCode},
		ExtraFields: schemas.BifrostErrorExtraFields{
			RequestType: schemas.ListModelsRequest,
		},
	}

	if !shouldReturnEmptyListModelsResponse(err) {
		t.Fatalf("expected unsupported_operation list-models error to return an empty success response")
	}
}

func TestShouldReturnEmptyListModelsResponseForOtherErrors(t *testing.T) {
	errCode := "internal_error"
	err := &schemas.BifrostError{
		Error: &schemas.ErrorField{Code: &errCode},
		ExtraFields: schemas.BifrostErrorExtraFields{
			RequestType: schemas.ListModelsRequest,
		},
	}

	if shouldReturnEmptyListModelsResponse(err) {
		t.Fatalf("expected non-unsupported-operation errors to be preserved")
	}
}

func TestShouldSkipProviderScopedListModelsRequestWhenDisabledByAllowedRequests(t *testing.T) {
	customProviderConfig := &schemas.CustomProviderConfig{
		AllowedRequests: &schemas.AllowedRequests{ListModels: false},
	}

	if !shouldSkipListModelsRequest("openai", customProviderConfig) {
		t.Fatalf("expected provider-scoped list-models request to be skipped when allowed_requests disables it")
	}
}

func TestShouldNotSkipProviderScopedListModelsRequestWhenEnabledByAllowedRequests(t *testing.T) {
	customProviderConfig := &schemas.CustomProviderConfig{
		AllowedRequests: &schemas.AllowedRequests{ListModels: true},
	}

	if shouldSkipListModelsRequest("openai", customProviderConfig) {
		t.Fatalf("expected provider-scoped list-models request to be allowed when allowed_requests enables it")
	}
}

func TestBuildDisabledListModelsResponse(t *testing.T) {
	got := buildDisabledListModelsResponse()
	if got.Message != "The model_list request is disabled for this provider." {
		t.Fatalf("expected disabled list-models response message, got %q", got.Message)
	}
	if len(got.Data) != 0 {
		t.Fatalf("expected disabled list-models response data to remain empty, got %#v", got.Data)
	}
}

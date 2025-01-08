// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	core "github.com/cohere-ai/cohere-go/v2/core"
	finetuning "github.com/cohere-ai/cohere-go/v2/finetuning"
	time "time"
)

type FinetuningListEventsRequest struct {
	// Maximum number of results to be returned by the server. If 0, defaults to 50.
	PageSize *int `json:"-" url:"page_size,omitempty"`
	// Request a specific page of the list results.
	PageToken *string `json:"-" url:"page_token,omitempty"`
	// Comma separated list of fields. For example: "created_at,name". The default
	// sorting order is ascending. To specify descending order for a field, append
	// " desc" to the field name. For example: "created_at desc,name".
	//
	// Supported sorting fields:
	//
	// - created_at (default)
	OrderBy *string `json:"-" url:"order_by,omitempty"`
}

type FinetuningListFinetunedModelsRequest struct {
	// Maximum number of results to be returned by the server. If 0, defaults to 50.
	PageSize *int `json:"-" url:"page_size,omitempty"`
	// Request a specific page of the list results.
	PageToken *string `json:"-" url:"page_token,omitempty"`
	// Comma separated list of fields. For example: "created_at,name". The default
	// sorting order is ascending. To specify descending order for a field, append
	// " desc" to the field name. For example: "created_at desc,name".
	//
	// Supported sorting fields:
	//
	// - created_at (default)
	OrderBy *string `json:"-" url:"order_by,omitempty"`
}

type FinetuningListTrainingStepMetricsRequest struct {
	// Maximum number of results to be returned by the server. If 0, defaults to 50.
	PageSize *int `json:"-" url:"page_size,omitempty"`
	// Request a specific page of the list results.
	PageToken *string `json:"-" url:"page_token,omitempty"`
}

type FinetuningUpdateFinetunedModelRequest struct {
	// FinetunedModel name (e.g. `foobar`).
	Name string `json:"name" url:"-"`
	// User ID of the creator.
	CreatorId *string `json:"creator_id,omitempty" url:"-"`
	// Organization ID.
	OrganizationId *string `json:"organization_id,omitempty" url:"-"`
	// FinetunedModel settings such as dataset, hyperparameters...
	Settings *finetuning.Settings `json:"settings,omitempty" url:"-"`
	// Current stage in the life-cycle of the fine-tuned model.
	Status *finetuning.Status `json:"status,omitempty" url:"-"`
	// Creation timestamp.
	CreatedAt *time.Time `json:"created_at,omitempty" url:"-"`
	// Latest update timestamp.
	UpdatedAt *time.Time `json:"updated_at,omitempty" url:"-"`
	// Timestamp for the completed fine-tuning.
	CompletedAt *time.Time `json:"completed_at,omitempty" url:"-"`
	// Timestamp for the latest request to this fine-tuned model.
	LastUsed *time.Time `json:"last_used,omitempty" url:"-"`
}

func (f *FinetuningUpdateFinetunedModelRequest) UnmarshalJSON(data []byte) error {
	type unmarshaler FinetuningUpdateFinetunedModelRequest
	var body unmarshaler
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	*f = FinetuningUpdateFinetunedModelRequest(body)
	return nil
}

func (f *FinetuningUpdateFinetunedModelRequest) MarshalJSON() ([]byte, error) {
	type embed FinetuningUpdateFinetunedModelRequest
	var marshaler = struct {
		embed
		CreatedAt   *core.DateTime `json:"created_at,omitempty"`
		UpdatedAt   *core.DateTime `json:"updated_at,omitempty"`
		CompletedAt *core.DateTime `json:"completed_at,omitempty"`
		LastUsed    *core.DateTime `json:"last_used,omitempty"`
	}{
		embed:       embed(*f),
		CreatedAt:   core.NewOptionalDateTime(f.CreatedAt),
		UpdatedAt:   core.NewOptionalDateTime(f.UpdatedAt),
		CompletedAt: core.NewOptionalDateTime(f.CompletedAt),
		LastUsed:    core.NewOptionalDateTime(f.LastUsed),
	}
	return json.Marshal(marshaler)
}

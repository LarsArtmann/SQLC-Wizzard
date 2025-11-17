package events

import (
	"time"
)

// ConfigurationUpdatedEventData represents typed event data for configuration updates
// Replaces 'any' type with specific typed structure
type ConfigurationUpdatedEventData struct {
	ProjectID       string         `json:"project_id"`
	ConfigurationID string         `json:"configuration_id"`
	UpdatedFields   []string       `json:"updated_fields"`
	OldValues       map[string]any `json:"old_values"`
	NewValues       map[string]any `json:"new_values"`
	UpdateReason    string         `json:"update_reason"`
	UpdatedBy       string         `json:"updated_by,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// EventType returns specific event type
func (cued *ConfigurationUpdatedEventData) EventType() string {
	return "configuration.updated"
}

// OccurredAt returns when configuration was updated
func (cued *ConfigurationUpdatedEventData) OccurredAt() time.Time {
	return cued.UpdatedAt
}

// AggregateID returns the configuration ID as the aggregate identifier
func (cued *ConfigurationUpdatedEventData) AggregateID() string {
	return cued.ConfigurationID
}

// Data returns the event data as a safe map for JSON serialization
func (cued *ConfigurationUpdatedEventData) Data() map[string]any {
	return map[string]any{
		"project_id":       cued.ProjectID,
		"configuration_id": cued.ConfigurationID,
		"updated_fields":   cued.UpdatedFields,
		"old_values":       cued.OldValues,
		"new_values":       cued.NewValues,
		"update_reason":    cued.UpdateReason,
		"updated_by":       cued.UpdatedBy,
		"updated_at":       cued.UpdatedAt,
	}
}

// Validate validates the configuration updated event data
func (cued *ConfigurationUpdatedEventData) Validate() error {
	if cued.ProjectID == "" {
		return &EventValidationError{
			Code:    "MISSING_PROJECT_ID",
			Message: "Project ID is required for configuration updated event",
		}
	}

	if cued.ConfigurationID == "" {
		return &EventValidationError{
			Code:    "MISSING_CONFIGURATION_ID",
			Message: "Configuration ID is required for configuration updated event",
		}
	}

	if len(cued.UpdatedFields) == 0 {
		return &EventValidationError{
			Code:    "NO_UPDATED_FIELDS",
			Message: "At least one updated field is required",
		}
	}

	if cued.OldValues == nil {
		return &EventValidationError{
			Code:    "MISSING_OLD_VALUES",
			Message: "Old values map is required",
		}
	}

	if cued.NewValues == nil {
		return &EventValidationError{
			Code:    "MISSING_NEW_VALUES",
			Message: "New values map is required",
		}
	}

	if cued.UpdateReason == "" {
		return &EventValidationError{
			Code:    "MISSING_UPDATE_REASON",
			Message: "Update reason is required for configuration updated event",
		}
	}

	if cued.UpdatedAt.IsZero() {
		return &EventValidationError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at time is required for configuration updated event",
		}
	}

	return nil
}

// NewConfigurationUpdatedEventData creates a new configuration updated event data with validation
func NewConfigurationUpdatedEventData(projectID, configurationID, updateReason, updatedBy string, updatedFields []string, oldValues, newValues map[string]any) (*ConfigurationUpdatedEventData, error) {
	if len(updatedFields) == 0 {
		return nil, &EventValidationError{
			Code:    "EMPTY_UPDATED_FIELDS",
			Message: "Updated fields list cannot be empty",
		}
	}

	if oldValues == nil {
		return nil, &EventValidationError{
			Code:    "NIL_OLD_VALUES",
			Message: "Old values map cannot be nil",
		}
	}

	if newValues == nil {
		return nil, &EventValidationError{
			Code:    "NIL_NEW_VALUES",
			Message: "New values map cannot be nil",
		}
	}

	// Validate that all updated fields exist in maps
	for _, field := range updatedFields {
		if _, exists := oldValues[field]; !exists {
			return nil, &EventValidationError{
				Code:    "FIELD_NOT_IN_OLD_VALUES",
				Message: "Field '" + field + "' not found in old values",
			}
		}
		if _, exists := newValues[field]; !exists {
			return nil, &EventValidationError{
				Code:    "FIELD_NOT_IN_NEW_VALUES",
				Message: "Field '" + field + "' not found in new values",
			}
		}
	}

	event := &ConfigurationUpdatedEventData{
		ProjectID:       projectID,
		ConfigurationID: configurationID,
		UpdatedFields:   updatedFields,
		OldValues:       oldValues,
		NewValues:       newValues,
		UpdateReason:    updateReason,
		UpdatedBy:       updatedBy,
		UpdatedAt:       time.Now(),
	}

	if err := event.Validate(); err != nil {
		return nil, err
	}

	return event, nil
}

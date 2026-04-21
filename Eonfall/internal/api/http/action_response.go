package http

import (
	"time"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

type ActionResponse struct {
	ID              uuid.UUID         `json:"id"`
	WorldID         uuid.UUID         `json:"world_id"`
	CivilizationID  uuid.UUID         `json:"civilization_id"`
	UserID          uuid.UUID         `json:"user_id"`
	ActionType      world.ActionType  `json:"action_type"`
	State           world.ActionState `json:"state"`
	TargetTick      int64             `json:"target_tick"`
	Payload         map[string]any    `json:"payload"`
	RejectionReason *string           `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	AppliedAt       *time.Time        `json:"applied_at,omitempty"`
}

func toActionResponse(a *world.WorldAction) ActionResponse {
	var payload map[string]any
	if a.Payload != nil {
		payload = a.Payload
	} else {
		payload = map[string]any{}
	}

	return ActionResponse{
		ID:              a.ID,
		WorldID:         a.WorldID,
		CivilizationID:  a.CivilizationID,
		UserID:          a.UserID,
		ActionType:      a.ActionType,
		State:           a.State,
		TargetTick:      a.TargetTick,
		Payload:         payload,
		RejectionReason: a.RejectionReason,
		CreatedAt:       a.CreatedAt,
		AppliedAt:       a.AppliedAt,
	}
}

type CivilizationContextResponse struct {
	DominantAxis  world.CivilizationAxis `json:"dominant_axis"`
	SecondaryAxis world.CivilizationAxis `json:"secondary_axis"`
	ProfileLabel  string                 `json:"profile_label"`
}

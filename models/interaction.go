package models

import "github.com/gocql/gocql"

type UserInteraction struct {
	UserID     gocql.UUID `json:"user_id"`
	ItemID     gocql.UUID `json:"item_id"`
	ActionType string     `json:"action_type"` // Ex: "view", "purchase", "rate"
}

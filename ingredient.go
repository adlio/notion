package notion

import "time"

type IngredientReport struct {
	IngredientID string    `json:"ingredient_id,omitempty"`
	Value        float64   `json:"value"`
	Date         time.Time `json:"date"`
}

type BatchIngredientReport struct {
	IngredientID string             `json:"ingredient_id"`
	Reports      []IngredientReport `json:"reports"`
}

type IngredientReportResponse struct {
	Errors []string `json:"errors"`
	Status string   `json:"status"`
}

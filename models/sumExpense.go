package models

type SumExpense struct {
	Sum      float64 `json:"sum"`
	Category string  `json:"category"`
}

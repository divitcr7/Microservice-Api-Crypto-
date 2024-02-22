package domain

type Symbol struct {
	Symbol  string `json:"symbol" db:"symbol"`
	Unicode rune   `json:"unicode" db:"unicode"`
}

package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/streamdp/ccd/repos"
)

// Symbols - validate the field so that the value is from the list of currencies
func Symbols(sr *repos.SymbolRepo) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		return sr.IsPresent(fl.Field().String())
	}
}

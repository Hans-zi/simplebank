package api

import (
	"github.com/Hans-zi/simple_bank/util"
	validator2 "github.com/go-playground/validator/v10"
)

var validCurrency validator2.Func = func(fl validator2.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupported(currency)
	}
	return false
}

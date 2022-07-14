package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsValidCurrency(currency)
	}
	return false
}

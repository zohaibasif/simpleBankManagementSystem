package util

const (
	EUR   = "EUR"
	USD   = "USD"
	RUPEE = "RUPEE"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case EUR, USD, RUPEE:
		return true
	}
	return false
}

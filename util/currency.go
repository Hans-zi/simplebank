package util

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	JPY = "JPY"
	KRW = "KRW"
)

func IsSupported(currency string) bool {
	switch currency {
	case USD, GBP, JPY, KRW, EUR:
		return true
	}
	return false
}

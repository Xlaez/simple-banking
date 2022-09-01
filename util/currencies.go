package util


	const NGN = "NGN"
	const USD = "USD"
	const EUR = "EUR"


func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, NGN, USD:
		return true;	
	}
	return false;
}
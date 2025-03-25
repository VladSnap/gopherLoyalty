package domain

type CurrencyUnit int

func (m CurrencyUnit) ToMajorUnit() float64 {
	return float64(m) / 100
}

func CurrencyFromMajorUnit(amount float64) CurrencyUnit {
	return CurrencyUnit(amount * 100)
}

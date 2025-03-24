package domain

type CurrencyUnit int

func (m CurrencyUnit) ToMajorUnit() float32 {
	return float32(m) / 100
}

func CurrencyFromMajorUnit(amount float32) CurrencyUnit {
	return CurrencyUnit(amount * 100)
}

package helpers

func GetOrDefaultInt(value *int, defaultValue int) int {
	if value == nil {
		return defaultValue
	}
	return *value
}

func GetOrDefaultFloat64(value *float64, defaultValue float64) float64 {
	if value == nil {
		return defaultValue
	}
	return *value
}

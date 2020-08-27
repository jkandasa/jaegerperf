package util

// PercentOf returns percentage
func PercentOf(value int, total int) float64 {
	return (float64(value) * float64(100)) / float64(total)
}

package utils

import "math"

// CalculateTotalPages calculates the total number of pages based on total items and limit per page
func CalculateTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

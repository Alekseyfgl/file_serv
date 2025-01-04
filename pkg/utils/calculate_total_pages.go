package utils

// CalculateTotalPages - calculates the total number of pages based on total count and page size
func CalculateTotalPages(totalCount, pageSize int) int {
	// If pageSize is less than or equal to zero, return 0 to avoid division by zero
	if pageSize <= 0 {
		return 0
	}
	// Use the formula to calculate the total number of pages, rounding up for any remaining items
	return (totalCount + pageSize - 1) / pageSize
}

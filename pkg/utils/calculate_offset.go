package utils

// CalculateOffset computes the number of items to skip based on the given page number and page size.
//
// Parameters:
//   - pageNumber: The current page number (1-based index).
//   - pageSize: The number of items per page.
//
// Returns:
//   - The offset (number of items to skip) for pagination.
//
// Example:
//
//	CalculateOffset(2, 10) returns 10 (skipping the first 10 items for page 2).
func CalculateOffset(pageNumber, pageSize int) int {
	return (pageNumber - 1) * pageSize
}

package utils

func GetPreviousPage(currentPage int) int {
	previousPage := currentPage - 1
	if previousPage < 1 {
		previousPage = 1
	}
	return previousPage
}

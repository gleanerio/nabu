package prune

// findMissingElements from chatGPT
func findMissingElements(a, b []string) []string {
	// Create a map to store the elements of ga.
	gaMap := make(map[string]bool)
	for _, s := range b {
		gaMap[s] = true
	}

	// Iterate through a and add any elements that are not in b to the result slice.
	var result []string
	for _, s := range a {
		if !gaMap[s] {
			result = append(result, s)
		}
	}

	return result
}

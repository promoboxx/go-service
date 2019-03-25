package unique

// Int32Slice returns a unique subset of the int32 slice provided.
func Int32Slice(input []int32) []int32 {
	result := make([]int32, 0, len(input))
	collector := make(map[int32]bool)

	for _, val := range input {
		if _, inSlice := collector[val]; !inSlice {
			collector[val] = true
			result = append(result, val)
		}
	}

	return result
}

// StringSlice returns a unique subset of the string slice provided
func StringSlice(input []string) []string {
	result := make([]string, 0, len(input))
	collector := make(map[string]bool)

	for _, val := range input {
		if _, inSlice := collector[val]; !inSlice {
			collector[val] = true
			result = append(result, val)
		}
	}

	return result
}

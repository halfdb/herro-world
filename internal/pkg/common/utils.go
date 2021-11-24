package common

func UniqueInt(input []int) []int {
	shown := make(map[int]bool)
	ints := make([]int, 0, len(input))
	for _, num := range input {
		if !shown[num] {
			shown[num] = true
			ints = append(ints, num)
		}
	}
	return ints
}

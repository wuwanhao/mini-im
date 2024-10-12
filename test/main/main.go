package main

// 泛型
func SumOfIntOrFloat[K comparable, V int | float64](m map[K]V) V {
	var result V
	for _, v := range m {
		result += v
	}
	return result
}

func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

package gopolitical

func PositiveModulus(a, b int) int {
	return (a%b + b) % b
}

package hashing

func BabyHash(in int) int {
	div := 3
	for i := 0; i < 3; i++ {
		in = in / div
		div = div + in
	}

	return div
}

package main

import (
	"debugging/hashing"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in, _ := strconv.Atoi(os.Args[1])
	out := hashing.BabyHash(in)

	fmt.Println(out)
}

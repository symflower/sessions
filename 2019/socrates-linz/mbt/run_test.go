package mbt

import (
	"os"
	"strconv"
	"testing"
)

func TestModelBasedTesting(t *testing.T) {
	iterations := 10
	if i := os.Getenv("ITERATIONS"); i != "" {
		ii, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}

		iterations = ii
	}

	Iterate(t, iterations)
}

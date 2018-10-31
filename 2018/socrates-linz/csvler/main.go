package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func process(data string) (int, error) {
	records, err := csv.NewReader(strings.NewReader(data)).ReadAll()
	if err != nil {
		return 0, err
	}

	overallSum := 0

	for i, row := range records {
		// Ignore the column row.
		if i == 0 {
			continue
		}

		sum := 0

		for _, column := range row {
			c, err := strconv.Atoi(column)
			if err != nil {
				return 0, err
			}

			sum += c
		}

		if sum > 10 {
			panic(fmt.Sprintf("Wups! The sum is too big: %d > 10", sum))
		}

		overallSum += sum
	}

	return overallSum, nil
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	sum, err := process(string(data))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sum: %d\n", sum)
}

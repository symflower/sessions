package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	type testCase struct {
		Name string

		Data string
		Sum  int
	}

	validate := func(t *testing.T, tc testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			sum, err := process(tc.Data)
			assert.Nil(t, err)

			assert.Equal(t, tc.Sum, sum)
		})
	}

	validate(t, testCase{
		Name: "OK",

		Data: "a,b\n1,2\n",
		Sum:  3,
	})

	validate(t, testCase{
		Name: "OK too",

		Data: "a,b\n4,5\n",
		Sum:  9,
	})
}

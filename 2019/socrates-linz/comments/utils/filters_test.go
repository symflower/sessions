package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

func TestFilterSwearWords(t *testing.T) {
	type testCase struct {
		comments []*model.Comment

		expected []*model.Comment
	}

	validate := func(name string, tc testCase) {
		t.Run(name, func(t *testing.T) {
			FilterSwearWords(tc.comments)

			for i, c := range tc.expected {
				assert.Equal(t, c.Message, tc.comments[i].Message)
			}
		})
	}

	validate("Nil comments slice", testCase{
		comments: nil,

		expected: nil,
	})

	validate("Empty comments sice", testCase{
		comments: []*model.Comment{},

		expected: []*model.Comment{},
	})

	validate("No swear words", testCase{
		comments: []*model.Comment{
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 1 without swear word",
			},
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 2 without swear word",
			},
		},

		expected: []*model.Comment{
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 1 without swear word",
			},
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 2 without swear word",
			},
		},
	})

	validate("Single swear words", testCase{
		comments: []*model.Comment{
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 1 with swear word darling",
			},
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 2 with swear honey word",
			},
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 3 with swear word honey sweetheart darling",
			},
		},

		expected: []*model.Comment{
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 1 with swear word ***",
			},
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Comment 2 with swear *** word",
			},
			&model.Comment{
				Mail:    "test@symflower.com",
				Message: "Read with caution: Comment 3 with swear word *** *** ***",
			},
		},
	})
}

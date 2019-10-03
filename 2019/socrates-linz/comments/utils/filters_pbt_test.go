package utils

import (
	"reflect"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/arbitrary"
	"github.com/leanovate/gopter/gen"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

func TestFilterSwearWordsPBT(t *testing.T) {
	type testCase struct {
		generator gopter.Gen
		condition interface{}

		parameters *gopter.TestParameters
	}

	validate := func(name string, tc testCase) {
		t.Run(name, func(t *testing.T) {
			if tc.parameters == nil {
				tc.parameters = gopter.DefaultTestParameters()
				tc.parameters.MinSuccessfulTests = 10
				tc.parameters.MaxSize = 20 // e.g. how many entries are at most in a slice
				tc.parameters.MinSize = 0  // e.g. how many entries are at least in a slice
			}
			properties := gopter.NewProperties(tc.parameters)

			arbitraries := arbitrary.DefaultArbitraries()
			if tc.generator != nil {
				arbitraries.RegisterGen(tc.generator)
			}
			properties.Property(name, arbitraries.ForAll(tc.condition))

			properties.TestingRun(t)
		})
	}

	validate("Random", testCase{
		condition: func(comments []*model.Comment) bool {
			t.Logf("Testing with comment slize of size %d\n", len(comments))
			for i, c := range comments {
				t.Logf("Message %d: %s\n", i, c.Message)
			}

			FilterSwearWords(comments)

			for _, c := range comments {
				for _, w := range swearwords {
					if strings.Contains(c.Message, w) {
						return false
					}
				}
			}

			return true
		},
	})

	validate("RandomAlphaString", testCase{
		condition: func(comments []*model.Comment) bool {
			t.Logf("Testing with comment slize of size %d\n", len(comments))
			for i, c := range comments {
				t.Logf("Message %d: %s\n", i, c.Message)
			}

			FilterSwearWords(comments)

			for _, c := range comments {
				for _, w := range swearwords {
					if strings.Contains(c.Message, w) {
						return false
					}
				}
			}

			return true
		},
		generator: gen.AlphaString(),
	})

	validate("Comments with swear words are filtered", testCase{
		generator: gen.StructPtr(reflect.TypeOf(&model.Comment{}), map[string]gopter.Gen{
			"Mail":    gen.Const("test@symflower.com"),
			"Message": gen.RegexMatch("^(.*(" + strings.Join(swearwords, "|") + ")?)*$"),
			"Created": gen.AnyTime(),
		}),
		condition: func(comments []*model.Comment) bool {
			t.Logf("Testing with comment slize of size %d\n", len(comments))
			for i, c := range comments {
				t.Logf("Message %d: %s\n", i, c.Message)
			}

			FilterSwearWords(comments)

			for _, c := range comments {
				for _, w := range swearwords {
					if strings.Contains(c.Message, w) {
						return false
					}
				}
			}

			return true
		},
	})

	validate("Comments with three to ten swear words are prefixed", testCase{
		condition: func(comments []*model.Comment) bool {
			t.Logf("Testing with comment slize of size %d\n", len(comments))
			for i, c := range comments {
				t.Logf("Message %d: %s\n", i, c.Message)
			}

			FilterSwearWords(comments)

			for _, c := range comments {
				if !strings.HasPrefix(c.Message, "Read with caution:") {
					return false
				}
			}

			return true
		},
		generator: gen.StructPtr(reflect.TypeOf(&model.Comment{}), map[string]gopter.Gen{
			"Mail":    gen.Const("test@symflower.com"),
			"Message": gen.RegexMatch("^(.*(" + strings.Join(swearwords, "|") + "){1}.*){3,10}$"),
			"Created": gen.AnyTime(),
		}),
	})

	validate("Comments with more than ten swear words are removed", testCase{
		condition: func(comments []*model.Comment) bool {
			t.Logf("Testing with comment slize of size %d\n", len(comments))
			for i, c := range comments {
				t.Logf("Message %d: %s\n", i, c.Message)
			}

			FilterSwearWords(comments)

			for _, c := range comments {
				if c.Message != "This message has been removed because it is to obscene." {
					return false
				}
			}

			return true
		},
		generator: gen.StructPtr(reflect.TypeOf(&model.Comment{}), map[string]gopter.Gen{
			"Mail":    gen.Const("test@symflower.com"),
			"Message": gen.RegexMatch("^(.*(" + strings.Join(swearwords, "|") + "){1}.*){11,}$"),
			"Created": gen.AnyTime(),
		}),
	})
}

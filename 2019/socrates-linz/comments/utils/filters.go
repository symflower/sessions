package utils

import (
	"strings"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

var swearwords = []string{"darling", "honey", "sweetheart"}

// FilterSwearWords replaces all swear words with "***" instead.
func FilterSwearWords(comments []*model.Comment) {
	for i, c := range comments {
		count := 0
		for _, s := range swearwords {
			count += strings.Count(c.Message, s)
			comments[i].Message = strings.ReplaceAll(c.Message, s, "***")
		}

		if count >= 3 && count <= 10 {
			comments[i].Message = "Read with caution: " + comments[i].Message
		} else if count > 10 {
			comments[i].Message = "This message has been removed because it is to obscene."
		}
	}
}

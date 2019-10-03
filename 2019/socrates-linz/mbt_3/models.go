package mbt

import (
	"fmt"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

func UserAdd(ctx *Context, user *model.User) {
	ctx.Users = append(ctx.Users, user)
}

func UserMailValid(ctx *Context) string {
	mails := make(map[string]bool, len(ctx.Users))
	for _, u := range ctx.Users {
		mails[u.Mail] = true
	}

	for {
		mail := fmt.Sprintf("user%d@symflower.com", ctx.Rand.Int())

		if _, ok := mails[mail]; !ok {
			return mail
		}
	}
}

func UserPasswordValid(ctx *Context) string {
	return "secret"
}

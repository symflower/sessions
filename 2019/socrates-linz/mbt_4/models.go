package mbt

import (
	"errors"
	"fmt"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

func UserAdd(ctx *Context, user *model.User) {
	ctx.Users = append(ctx.Users, user)
}

var ErrUserMailExists = errors.New("does already exist")

var userMailChoices = []func(ctx *Context) (string, error, bool){
	func(ctx *Context) (string, error, bool) {
		return UserMailValid(ctx), nil, true
	},
	func(ctx *Context) (string, error, bool) {
		mail, ok := UserMailExists(ctx)
		if !ok {
			return "", nil, false
		}

		return mail, ErrUserMailExists, true
	},
}

func UserMail(ctx *Context) (string, error) {
	for {
		c := userMailChoices[ctx.Rand.Int()%len(userMailChoices)]

		mail, err, ok := c(ctx)
		if !ok {
			continue
		}

		return mail, err
	}
}

func UserMailExists(ctx *Context) (string, bool) {
	if len(ctx.Users) == 0 {
		return "", false
	}

	user := ctx.Users[ctx.Rand.Int()%len(ctx.Users)]

	return user.Mail, true
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

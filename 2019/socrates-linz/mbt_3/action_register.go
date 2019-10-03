package mbt

import (
	"strings"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

func init() {
	ActionRegister(&Action{
		Name: "RegisterGet",
		Call: func(ctx *Context) {
			RegisterGet(ctx)
		},
	})
}

func RegisterGet(ctx *Context) {
	ctx.Write("mbt.RegisterGet(ctx)\n")

	_, body := HTTPGetValid(ctx, "/register")

	dom := DOM(ctx, body)
	if dom.Find("#form_register") == nil {
		ctx.Fatal("Register form does not exist")
	}
}

func init() {
	ActionRegister(&Action{
		Name: "RegisterPost",
		Call: func(ctx *Context) {
			HTTPPostDataSet(ctx, "mail", UserMailValid(ctx))
			HTTPPostDataSet(ctx, "password", UserPasswordValid(ctx))

			RegisterPost(ctx)
		},
	})
}

func RegisterPost(ctx *Context) {
	ctx.Write("mbt.RegisterPost(ctx)\n")

	form := ctx.FormData

	_, body := HTTPPostSend(ctx, "/register")

	if !strings.Contains(body, "You registered the user "+form["mail"]) {
		ctx.Fatal("Registered message does not exist")
	}

	ctx.Users = append(ctx.Users, &model.User{
		Mail:     form["mail"],
		Password: form["password"],
	})
}

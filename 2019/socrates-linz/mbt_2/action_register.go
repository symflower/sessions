package mbt

import (
	"strings"
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
			HTTPPostDataSet(ctx, "mail", "user@symflower.com")
			HTTPPostDataSet(ctx, "password", "secret")

			RegisterPost(ctx)
		},
	})
}

func RegisterPost(ctx *Context) {
	ctx.Write("mbt.RegisterPost(ctx)\n")

	_, body := HTTPPostSend(ctx, "/register")

	if !strings.Contains(body, "You registered the user user@symflower.com") {
		ctx.Fatal("Registered message does not exist")
	}
}

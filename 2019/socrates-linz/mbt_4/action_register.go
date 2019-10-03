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
			mail, mailError := UserMail(ctx)
			HTTPPostDataSet(ctx, "mail", mail)
			if mailError != nil {
				HTTPPostErrorSet(ctx, "mail", mailError.Error())
			}

			HTTPPostDataSet(ctx, "password", UserPasswordValid(ctx))

			RegisterPost(ctx)
		},
	})
}

func RegisterPost(ctx *Context) {
	ctx.Write("mbt.RegisterPost(ctx)\n")

	form := ctx.FormData
	formErrors := ctx.FormErrors

	_, body := HTTPPostSend(ctx, "/register")

	dom := DOM(ctx, body)

	if mailError, ok := formErrors["mail"]; ok {
		m := dom.Find("#mail_error")
		if m == nil || !strings.Contains(m.Text(), mailError) {
			ctx.Fatalf("Cannot find error message %q", mailError)
		}
	}

	if len(formErrors) > 0 {
		if dom.Find("#form_register") == nil {
			ctx.Fatal("Register form does not exist")
		}
	} else {
		if !strings.Contains(body, "You registered the user "+form["mail"]) {
			ctx.Fatal("Registered message does not exist")
		}

		ctx.Users = append(ctx.Users, &model.User{
			Mail:     form["mail"],
			Password: form["password"],
		})
	}
}

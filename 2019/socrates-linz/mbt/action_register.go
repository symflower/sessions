package mbt

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

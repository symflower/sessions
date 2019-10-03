package mbt

func init() {
	ActionRegister(&Action{
		Name: "IndexGet",
		Call: func(ctx *Context) {
			IndexGet(ctx)
		},
	})
}

func IndexGet(ctx *Context) {
	ctx.Write("mbt.IndexGet(ctx)\n")

	_, body := HTTPGetValid(ctx, "/")

	dom := DOM(ctx, body)
	if dom.Find("#form_comment") == nil {
		ctx.Fatal("Comment form does not exist")
	}
}

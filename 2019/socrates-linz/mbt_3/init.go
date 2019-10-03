package mbt

import (
	"time"

	"github.com/symflower/sessions/2019/socrates-linz/comments/server"
)

func Init(ctx *Context) {
	// Setup the server.
	go func() { // Run the server in the background.
		server.Run()
	}()

	time.Sleep(time.Second) // Wait a second to let the server start.
}

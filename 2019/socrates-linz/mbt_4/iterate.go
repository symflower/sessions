package mbt

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

type Action struct {
	Name string
	Call func(ctx *Context)
}

var actions []*Action

func ActionRegister(action *Action) {
	for _, a := range actions {
		if a.Name == action.Name {
			panic(fmt.Sprintf("action %q does already exist", action.Name))
		}
	}

	actions = append(actions, action)
}

type Context struct {
	t        *testing.T
	testFile *os.File

	Rand *rand.Rand

	Body string

	FormData   map[string]string
	FormErrors map[string]string
	Users      []*model.User
}

func NewContext(t *testing.T) *Context {
	return &Context{
		t: t,

		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),

		FormData:   map[string]string{},
		FormErrors: map[string]string{},
	}
}

func (ctx *Context) TestFileCreate() {
	f, err := ioutil.TempFile("", "mbt_*_test.go")
	if err != nil {
		panic(err)
	}
	ctx.testFile = f

	fmt.Fprintf(os.Stderr, "Write test case to %s\n", f.Name()) // Write to STDERR because the testing log is not shown on errors.

	ctx.Write(`package main_test

import (
	"testing"

	"github.com/symflower/sessions/2019/socrates-linz/mbt"
)

func TestNewCase(t *testing.T) {
	ctx := mbt.NewContext(t)

	mbt.Init(ctx)

	` + "\n")
}

func (ctx *Context) Exit() {
	if ctx.testFile != nil {
		ctx.Write("}\n")
		ctx.testFile.Close()
	}
}

func (ctx *Context) Write(data string) {
	fmt.Print(data)

	if ctx.testFile != nil {
		ctx.testFile.WriteString(data)
	}
}

func (ctx *Context) Writef(format string, args ...interface{}) {
	ctx.Write(fmt.Sprintf(format, args...))
}

func (ctx *Context) Fatal(args ...interface{}) {
	ctx.Exit()

	fmt.Println(ctx.Body)

	ctx.t.Fatal(args...)
}

func (ctx *Context) Fatalf(format string, args ...interface{}) {
	ctx.Fatal(fmt.Sprintf(format, args...))
}

func Iterate(t *testing.T, iterations int) {
	if len(actions) == 0 {
		t.Fatal("no actions defined")
	}

	ctx := NewContext(t)
	defer ctx.Exit()

	Init(ctx)

	ctx.TestFileCreate()

	// Start the model-based testing iterations.
	for {
		action := actions[ctx.Rand.Int()%len(actions)]

		action.Call(ctx)

		if iterations > 0 {
			iterations--
		}
		if iterations == 0 {
			break
		}
	}
}

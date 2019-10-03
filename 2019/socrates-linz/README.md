# Let’s Tackle Software Testing (by Symflower)

The presentation for this repository/workshop can be found in [Symflower - Let's Tackle Software Testing at Socrates 2019.pdf](Symflower - Let's Tackle Software Testing at Socrates 2019.pdf). Please take a look for the introduction and additional details to the examples below and the repository in general.

# VERY IMPORTANT

This repository holds code that is only one thing: an example for analyzing software and testing software. Do not use any of the application's code or patterns in production. The whole application deliberately includes bugs, bad code quality and flaws.

# Installation

If you want to change the source code of the examples, which is not necessary to execute the examples, you need an editor or IDE that supports Go which usually also means that you need Go (we use the latest Go 1.12 for this repository) in your operating system available. Please have a look at https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins for common editor and IDE integrations, e.g. we use [Visual Studio Code](https://code.visualstudio.com/) with the [vscode-go plugin](https://github.com/Microsoft/vscode-go).

This repository comes with its own Docker file, which holds all dependencies and tools that are mentioned in the presentation and in the README. For installing Docker please see the installation guide of your operating system/distribution or https://docs.docker.com/install/ (e.g. https://docs.docker.com/install/linux/docker-ce/ubuntu/). You can either use the Docker image uploaded on https://hub.docker.com with the following command ...

```bash
docker pull symflower/socrates-linz-2019:latest
```

... or build the Docker image yourself using the following command:

```bash
docker build --tag symflower/socrates-linz-2019:latest .
```

Now we can use the Docker image to execute all examples of this repository. The following command starts a Docker container binding the current directory into the container as well as forwarding the port 8080 of the container to your host machine.

```bash
docker run --name socrates --rm -it --publish 8080:8080 --volume $PWD:/go/src/github.com/symflower/sessions/2019/socrates-linz symflower/socrates-linz-2019:latest bash
```

Now you can start the example web server with the following command and access it with your browser under http://localhost:8080/.

```bash
go run comments/main.go
```

# Examples

The sample application comes with unit tests, system tests as well as the base of a model based tester.

## Unit tests traditional and with Property Based Testing(PBT)

The application is enjoying great popularity. That is why our product manager has requested a new feature to make it also suitable for children. This feature should ensure that all swear words are filtered.

Before this functionality is integrated into the production code, we want to make sure it is working as intended. Take a look at file `comments/utils/filters.go` to checkout the current implementation. It comes with four traditional unit tests written in a table-driven format (see `comments/utils/filters_test.go`).

Table-driven tests have the advantage that no code is duplicated, they are easy to read and the data for testing is stored in one place. The structure you see in the file is our (Symflower's) style of writing table-driven tests: we have a test case struct and a "validate" function that does all the execution and checks for each test case. Doing one explicit call per test case, instead of a loop on a list of test cases, has the advantage that we have a useful stack trace on any error. Side note: we also group inputs and outputs of a test case.

Run the tests with the following command:
```bash
go test -v -run=TestFilterSwearWords$ -coverprofile=coverage.out -coverpkg=github.com/symflower/sessions/2019/socrates-linz/comments/... github.com/symflower/sessions/2019/socrates-linz/comments/utils
```

Generate a human readable report with the following command:
```bash
go tool cover -html=coverage.out -o coverage.html
```

The current test suite for our function does not reach full line coverage, because we are missing a comment with more than 10 swear words. Of course we could add this test case now and be done with testing. However, what if we add more features? What if we overlooked some test case? Creating the existing test cases was alread hard work. Let's try to recreate the same tests + the missing case but this time with automatically.

In traditional unit testing concrete inputs and outputs are defined. Finding all relevant ones is cumbersome, and often corner cases are missed. Property based testing follows the notion to no longer define concrete inputs and outputs. The tester rather defines "properties" that have to hold for certain inputs.

```
forAll (x uint[]) holds sum(x) > 0
```

The above property defines that for all unsigned integer slices holds, that the some of its elements needs to be greater than zero. A property based testing framework randomly generates slices of uints and checks that the defined property `sum(x)>0` holds. If all conditions hold, a valid test case has been found.

In our opinion property based testing cannot replace unit testing. Not all functions are suited for property based testing. But oftentimes it might be helpful as complimentary testing technique, revealing bugs that your unit tests missed. And it also helps in documenting the intended behavior of a function.

[QuickCheck](https://en.wikipedia.org/wiki/QuickCheck) is a well known property based testing framework for Haskell that exists for all of the popular programming languages. We will use [Gopter](https://github.com/leanovate/gopter) one of the implementations for Go.

Take a look at the function in `comments/utils/filters.go` and its accompanying tests in `comments/utils/filters_pbt_test.go`.

A test case consists of a `generator` for random values, a `condition` that is checked for each generated value and `parameters` that configures the model based tester, i.e. how many test cases, how many shrinking steps (to reduce found test cases to a minimum which leads to the same test case), the seed for random data generation.

The validate function starts a subtest for each defined test case.

The first test case `TestFilterSwearWords/Random` only specifies the condition and leaves the other test parameters to it's defaults. There is some logging output to inspect the generated values.

Run the following commands to check the generated outputs and inspect its coverage:

```bash
go test -v -run=TestFilterSwearWords/Random$ -coverprofile=coverage.out -coverpkg=github.com/symflower/sessions/2019/socrates-linz/comments/... github.com/symflower/sessions/2019/socrates-linz/comments/utils
```

Play around with the configuration in `tc.parameters` to check wether you can increase coverage by producing more tests or by increasing the size of generated inputs. When inspecting the generated outputs you might end up at the conclusion: with purely random values it is very unlikely to eventually cover the `if-else` block in `FilterSwearWords`.

Next run the test `TestFilterSwearWords/RandomAlphaString`. The outputs look a bit more promising, but still the likelyhood of reaching full coverage is rather small. This brings us to the last three examples that rely on a regular expression that finally increases the likelyhood of enough swear words to also cover the remaining blocks.

Which of the tests would you actually use? The problem I see with  `TestFilterSwearWords/Comments_with_swear_words_are_filtered` and  `TestFilterSwearWords/Comments_with_more_than_ten_swear_words_are_removed` is that they only add entries with many swear words. There are no entries with zero swear words.

One downside of property based testing is, that finding good properties is hard. Not all functions are completely suitable for property based testing since it is often necessary to almost introduce the same implementation again. For a lot of functionality adding properties can help guide the test generation process which is very similar to how fuzzing works in general. However, "mutation based" fuzzing still has the advantage that existing test cases are mutated which is easier to define, since no properties are needed, and the generation can learn from found test cases ("coverage guided fuzzing"). In our opinion, a combination of mutation based fuzzing, covearge guided fuzzing and property based testing is needed to be more useful in generating test cases automatically and easily.

## System tests and their code coverage

Let's look at the file "comments/main_test.go". As you can see the application is already automatically tested. All requirements of the specification are included in our test cases. We can register a user, a registered user can post a message and we can view comments. These test cases are based on querying a running web server that is initialized at the beginning of running the test cases.

Before we dig deeper, let's run all system test cases with the following command including an analysis which code is covered by these tests:

```bash
go test -v -run=TestApplication -coverprofile=coverage.out -coverpkg=github.com/symflower/sessions/2019/socrates-linz/comments/... github.com/symflower/sessions/2019/socrates-linz/comments
```

The file "coverage.out" now holds the raw coverage of all code that has been executed. We run the following command to create an HTML representation that can be opened with any browser.

```bash
go tool cover -html=coverage.out -o coverage.html
```

### Question: What observations can you make by looking at the coverage?

Some pointers:
- The coverage tool is "line based". Hence, we cannot be sure that a line is really fully covered, e.g. what about different conditions of an "if" statement?
- Even though we have all requirements of the specifications covered, there is still a lot of code not covered, e.g. pretty much all error cases are not covered.
- We do not see the coverage of the templates, so we do not know what was executed in the templates.
- Some statements have multiple outcomes, e.g. an SQL statement could fail, still the line is marked as covered.

Additional questions:
- Should we test more?
- Do we have code that has been covered but not tested?

## Mutation testing: Do we have code that has been covered but not tested?

With mutation testing we can answer the question how good the quality of our test suite is. Hence, we can answer the question if we have covered code in our test cases but not really tested that this code has been covered. This is immensly helpful in finding code that require test cases but also code that might be simply unneeded.

We will use the mutation testing tool [go-mutesting](https://github.com/zimmski/go-mutesting). You can find mutation testing tools for other languages at https://github.com/theofidry/awesome-mutation-testing. The following command evaluates our test suite. Please note, this command may take while to finish. The file "testdata/go-mutesting.log" holds an example output of the command (which depends on the repository content, so you might see some different output).

```bash
go-mutesting --verbose --exec scripts/exec-mutation.sh github.com/symflower/sessions/2019/socrates-linz/comments/...
```

The final information of the output is the report on the mutation score and how many mutations have been passed and failed. For example this could be your output, which depends on the content of the repository, "The mutation score is 0.203125 (13 passed, 47 failed, 14 duplicated, 4 skipped, total is 64)" which tells us that there where 13 mutations (passed) that were caught by our test suite, but 47 mutations (failed) that have not been caught.

When we compare these mutations with our coverage we see that lots of them apply to "not covered" lines of code. Hence, we would have guested the same missing test cases with a simple coverage tool that looks at e.g. the line coverage. However, some of these mutations belong to lines that are covered. Therefore, code that is executed by our tests but not checked. Let's look at some of these mutations to understand a.) how we could missed the mutations b.) how we can test them automatically in the future.

The following mutation tells us that we do not test a log statement. Question: Is this a problem?

```diff
--- /go/src/github.com/symflower/sessions/2019/socrates-linz/comments/main.go   2019-10-02 15:48:24.846246973 +0000
+++ /tmp/go-mutesting-145288687//go/src/github.com/symflower/sessions/2019/socrates-linz/comments/main.go.8     2019-10-02 16:09:17.136799841 +0000
@@ -40,8 +40,7 @@

        http.HandleFunc("/", middleware(controller.HandleIndex))
        http.HandleFunc("/register", middleware(controller.HandleRegister))
-
-       log.Printf("Listening on port 8080")
+       _ = log.Printf
        err = http.ListenAndServe(":8080", nil)
        if err != nil {
                panic(err)
```

The following mutation tells us that we do not test if there are any errors for the mail address nor password of the comment form. Is this a problem?

```diff
--- /go/src/github.com/symflower/sessions/2019/socrates-linz/comments/controller/index.go       2019-10-02 15:59:42.619636275 +0000
+++ /tmp/go-mutesting-145288687//go/src/github.com/symflower/sessions/2019/socrates-linz/comments/controller/index.go.10        2019-10-02 16:09:24.280814227 +0000
@@ -96,7 +96,7 @@
                data.Form.MessageError = "This is a required field."
        }

-       if data.Form.MailError == "" && data.Form.PasswordError == "" && data.Form.MessageError == "" {
+       if true && data.Form.MessageError == "" {
                err := model.CommentAdd(db, data.Form.Mail, data.Form.Message)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
```

The following mutation tells us that we do not test the reset code for the comment form. Is this a problem?

```diff
--- /go/src/github.com/symflower/sessions/2019/socrates-linz/comments/controller/index.go       2019-10-02 15:59:42.619636275 +0000
+++ /tmp/go-mutesting-145288687//go/src/github.com/symflower/sessions/2019/socrates-linz/comments/controller/index.go.25        2019-10-02 16:09:28.992823716 +0000
@@ -103,8 +103,7 @@

                        return
                }
-
-               data.Form.Mail = ""
+               _ = data.Form.Mail
                data.Form.Password = ""
                data.Form.Message = ""
        }
```

The following mutation tells us that we do not care about the value of the password. How can that be?

```diff
--- /go/src/github.com/symflower/sessions/2019/socrates-linz/comments/controller/index.go       2019-10-02 15:59:42.619636275 +0000
+++ /tmp/go-mutesting-145288687//go/src/github.com/symflower/sessions/2019/socrates-linz/comments/controller/index.go.15        2019-10-02 16:09:25.116815911 +0000
@@ -72,7 +72,7 @@
        data := templateIndexData{}

        data.Form.Mail = r.FormValue("mail")
-       data.Form.Password = r.FormValue("password")
+       _, _ = data.Form.Password, r.FormValue
        data.Form.Message = r.FormValue("message")

        if data.Form.Mail == "" {
```

Hint: Tests can have bugs too. Question: Can you spot the bug?

These mutations make it painfully clear: Mutation testing must be part of the testing process if automated tests are written manually or semi-automatically, or else we cannot make sure that we have enough test cases and we cannot make sure that our test cases are testing what we thought is tested.

### Digging into the system test code

Finally, let's have a deeper look at test system test code.

Question: What can we improve?

Question: How easy is it to add new test cases or even whole scenarios? What if we want to repeat some action (e.g. add 10 comments)? Are we bound to linear scenarios or can we branch out?

Question: When can we decided that we added enough scenarios? How can we come up with new test scenarios without wasting our time?

Question: How can we detect if one request changes the behavior of another request? In other words: How can we make sure that we have not overlooked a side effect?

Hint: Let's define a model to generate whole test scenarios without thinking about how these scenarios are defined. We need a model that has the potential to automatically find side effects.

## Model-based testing: How can we add tests that find problems semi-automatically?

The problem with the usualy definition of test scenarios is that humans have to keep a lot of context in their mind to create them, e.g. before one can buy an item in a shop, the item and user have to be created, the item must be put into the basket and the basket must be successfully checked out. This leads to a massiv question mark about all the other scenarios, e.g. what happens if one of the scenario phases is forgotten? Model-based testing solves this problem by keeping the focus only on one specific action of the whole application. One action is one possibility of the whole application, e.g. pressing the register button to create a new user.

Bascially, we are focusing on the actions that can be done, the data that needs to be generated to do this actions and the validations to check if the actions can be done at all and were done correctly.

Usually model-based testing tools can ease the creation and execution of such models, e.g. https://graphwalker.github.io/ and https://github.com/zimmski/tavor. However, since we only want to demonstrate how easy the creation of such a model-based testing suite is, we will do it on our own. Still, if we would use an existing tool/framework to do model-based testing, we would save a lot of time writing helpers and logic for the model-based testing functionality. Also, the API would be a given and we do not have to think about it. However, this is an example on how easy we can write such a tool ourselves.

### Base of our model-based tester: actions for the index and register pages via GET

We already prepared iterations of a manually created model-based testing tool that generates test scenarios for our application. Have a look at the "mbt/run_test.go" file. The file holds one test case that runs the model-based testing tool with 10 iterations. The following commands runs the test case while taking note of the coverage.

```bash
go test -v -run=TestModelBasedTesting -coverprofile=coverage.out -coverpkg=github.com/symflower/sessions/2019/socrates-linz/comments/... ./mbt
```

One output could be:

```
=== RUN   TestModelBasedTesting
2019/10/03 10:54:20 Listening on port 8080
Write test case to /tmp/mbt_325978464_test.go
package main_test

import (
        "testing"

        "github.com/symflower/sessions/2019/socrates-linz/mbt"
)

func TestNewCase(t *testing.T) {
        ctx := mbt.NewContext(t)

        mbt.Init(ctx)


mbt.IndexGet(ctx)
2019/10/03 10:54:21 Request / with GET
2019/10/03 10:54:21 Done
mbt.RegisterGet(ctx)
2019/10/03 10:54:21 Request /register with GET
2019/10/03 10:54:21 Done
mbt.IndexGet(ctx)
2019/10/03 10:54:21 Request / with GET
2019/10/03 10:54:21 Done
mbt.IndexGet(ctx)
2019/10/03 10:54:21 Request / with GET
2019/10/03 10:54:21 Done
mbt.RegisterGet(ctx)
2019/10/03 10:54:21 Request /register with GET
2019/10/03 10:54:21 Done
mbt.RegisterGet(ctx)
2019/10/03 10:54:21 Request /register with GET
2019/10/03 10:54:21 Done
mbt.IndexGet(ctx)
2019/10/03 10:54:21 Request / with GET
2019/10/03 10:54:21 Done
mbt.IndexGet(ctx)
2019/10/03 10:54:21 Request / with GET
2019/10/03 10:54:21 Done
mbt.IndexGet(ctx)
2019/10/03 10:54:21 Request / with GET
2019/10/03 10:54:21 Done
mbt.RegisterGet(ctx)
2019/10/03 10:54:21 Request /register with GET
2019/10/03 10:54:21 Done
}
--- PASS: TestModelBasedTesting (1.02s)
PASS
coverage: 49.6% of statements in github.com/symflower/sessions/2019/socrates-linz/comments/...
ok      github.com/symflower/sessions/2019/socrates-linz/mbt    1.019s  coverage: 49.6% of statements in github.com/symflower/sessions/2019/socrates-linz/comments/...
```

The log tells us that the tool has generated requests to the index and register pages via GET multiple times, leading to a coverage of 49.6%. Also, this test case has been written into the file "/tmp/mbt_325978464_test.go". Which holds the following content.

```go
package main_test

import (
        "testing"

        "github.com/symflower/sessions/2019/socrates-linz/mbt"
)

func TestNewCase(t *testing.T) {
        ctx := mbt.NewContext(t)

        mbt.Init(ctx)


mbt.IndexGet(ctx)
mbt.RegisterGet(ctx)
mbt.IndexGet(ctx)
mbt.IndexGet(ctx)
mbt.RegisterGet(ctx)
mbt.RegisterGet(ctx)
mbt.IndexGet(ctx)
mbt.IndexGet(ctx)
mbt.IndexGet(ctx)
mbt.RegisterGet(ctx)
}
```

This test case can be executed using the following commands:

```bash
mkdir run
mv /tmp/mbt_325978464_test.go run/
go test -v run/*.go
```

As you can see, the exact same scenario, without printing the code of the test case, has been executed. Therefore, our tool already generates random test scenarios and also writes executable test cases. Hence, we have already reached the goal that we do not have to think about scenarios.

One way of using this tool is to add new actions and checks, and then generate test case for hours and add every generated case to a repository when a new coverage is reached.

Now let's have a look at the source code of the tool. The file "mbt/run_test.go" tells us that we can use the environment variable "ITERATIONS" to set how many iterations we want the tool to run for one test case. Digging deeper into the function "Iterate" in the file "mbt/iterate.go" tells us that we can set "ITERATIONS" to "-1" to run indefinitelly instead of a fixed count. We also see in this function that the code takes a random "action" and executes it. What are these actions?

Our current actions are defined in "mbt/action_index.go" and "mbt/action_register.go". As you can see we register an action which prints out the source code to reproduce the action's call (i.e. "ctx.Write(...)") and we den proceed to do a HTTP get request that must be valid and search for a form in the response. Simply put, we browse our application with our tool but do not type in any data.

Since this is an example, we do not think too much about the structure of our tool, e.g. how does the code look if we have 1000 pages? Would we then add a 1000 files in the same directory? These are questions that have to be answered with another example. Right now, we only care about the basic functionality of our model-based tester.

Our current code already shows the most important attribute of model-based testing: We do not think about in which order the actions(in our case pages of the application) are called but only that they exist and how they can be reached. Since all of our pages can be directly reached we simply add them directly to our tester.

### Register a user

Next step: let's register a user. (The finished code can be found in mbt_2.)

Let's add the following code to our file "mbt/action_register.go".

```go
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
```

This code adds a new action "RegisterPost" to our set of actions. Hence, the action can be directly called at any moment in a test scenario. The action sets the mail and password data with valid values for a post request and then calls "RegisterPost" which executes the post request and validates that the registered message is shown. Let's execute our model-based tester with our new action a few times. You will run into the following problem.

```
=== RUN   TestModelBasedTesting
2019/10/03 11:52:41 Listening on port 8080
Write test case to /tmp/mbt_937246041_test.go
package main_test

import (
        "testing"

        "github.com/symflower/sessions/2019/socrates-linz/mbt"
)

func TestNewCase(t *testing.T) {
        ctx := mbt.NewContext(t)

        mbt.Init(ctx)


mbt.IndexGet(ctx)
2019/10/03 11:52:42 Request / with GET
2019/10/03 11:52:42 Done
mbt.IndexGet(ctx)
2019/10/03 11:52:42 Request / with GET
2019/10/03 11:52:42 Done
mbt.RegisterGet(ctx)
2019/10/03 11:52:42 Request /register with GET
2019/10/03 11:52:42 Done
mbt.RegisterGet(ctx)
2019/10/03 11:52:42 Request /register with GET
2019/10/03 11:52:42 Done
mbt.HTTPPostDataSet(ctx, "mail", "user@symflower.com")
mbt.HTTPPostDataSet(ctx, "password", "secret")
mbt.RegisterPost(ctx)
2019/10/03 11:52:42 Request /register with POST
2019/10/03 11:52:42 Done
mbt.RegisterGet(ctx)
2019/10/03 11:52:42 Request /register with GET
2019/10/03 11:52:42 Done
mbt.IndexGet(ctx)
2019/10/03 11:52:42 Request / with GET
2019/10/03 11:52:42 Done
mbt.HTTPPostDataSet(ctx, "mail", "user@symflower.com")
mbt.HTTPPostDataSet(ctx, "password", "secret")
mbt.RegisterPost(ctx)
2019/10/03 11:52:42 Request /register with POST
2019/10/03 11:52:42 Done
}

...

                        <form id="form_register" class="pure-form pure-form-aligned" method="POST">
                                <fieldset>
                                        <legend>Register a user or <a href="/">go commenting with your user</a></legend>

                                        <div class="pure-control-group">
                                                <label for="mail">Mail address</label>
                                                <input id="mail" name="mail" type="email" placeholder="Mail address" value="user@symflower.com">
                                                <span id="mail_error" class="pure-form-message-inline">The mail address user@symflower.com does already exist.</span>
                                        </div>

                                        <div class="pure-control-group">
                                                <label for="password">Password</label>
                                                <input id="password" name="password" type="password" placeholder="Password" value="secret">

                                        </div>

                                        <div class="pure-controls">
                                                <input type="submit" value="Create" class="pure-button pure-button-primary">
                                        </div>
                                </fieldset>
                        </form>

        </body>
</html>

}
--- FAIL: TestModelBasedTesting (1.02s)
    iterate.go:104: Registered message does not exist
FAIL
coverage: 62.3% of statements in github.com/symflower/sessions/2019/socrates-linz/comments/...
FAIL    github.com/symflower/sessions/2019/socrates-linz/mbt    1.024s
```

As you can see from the output the mail address "user@symflower.com" is already defined. The reason is simple: the model-based tester has tried to register the same user twice and we do not expect the registration of a user to fail. This is one important attribute of model-based testing: if some reaction of the system under test is unexpected a good model-based tester will always stop the execution while reporting the unexpected behavior.

### Register random users with unique mail addresses

Next step: let's register users with unique mail addresses. (The finished code can be found in mbt_3.)

Hardcoded data is usually a bad practice sign in model-based testing and in testing in general. Hence, we will add a generator for our mail addresses and since we need to remember which addresses were registered, we need to save registered addresses in our model-based tester state.

Let's change some files and then discuss the changes. First, adapt our "RegisterPost" code in "mbt/action_register.go" to the following.

```go
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
```

Additionally, add the file "mbt/models.go" with the following content.

```go
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
```

Instead of using static data in our action code we now use functions that "generate" data for us. The logic of the generation is basically hidden from the point of view of the action. This is important because the action should not care about how a mail address is valid, it only care that it gets a valid mail address. We also see that we now keep track of the added users. This should not be a complete copy of our application logic for saving a user. We only want to keep track of what is important to us during the application. Sinc we need a mail address and a password to log in, we need to keep both. The file "mbt/models.go" holds the handling of our data. Again, this is not an exercise about how to structure such keeping of data or how the project should be structured. The function "UserMailValid" implements the logic of generating a new valid mail address.

Question: Do you see a problem with this function? Hint: One problem is that we still a very static about what mail addresses are valid "user%d@symflower.com", e.g. does not test different domains.

Executing our new code of our model-based tester reveals that we now can register multiple users without any problems. Hence, if a bug would involve multiple users, we have now the propability to automatically detect it. However, we are still not done, e.g. what about handling registrations with existing mail addresses?

### Handling existing mail addresses during registration

Next step: let's handle the existing mail address error during the registration. (The finished code can be found in mbt_4.)

We adapt again some files. Change the code of our RegisterPost in "mbt/action_register.go" to the following.

```go
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
```

Add the following code to "mbt/models.go".

```go
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
```

The changes allow our action to not just use a valid mail address but instead get either a valid or an existing mail address. Additionally, we can also receive an error for a mail address which is then checked in our POST code. Simply put, if there is a form error, we check that the error exists in the output, and if there is no error, we successfully registered a user. The generation code is now slightly different. Since we either want a valid or an existing mail address we need to create a random choice (see "UserMail"). However, there is one big problem: what if no user has been registered yet? Then we simply abort the choice and choose again. Since we know that either a user must exist or not in our database of the application, we can make such a decision without corrupting our ability to fully test our application.

Exercise: What do we need to add to our source code so handle the error of invalid mail addresses in the registration? Hint: We only need to add an error and another choice for our mail generation. The rest of the model-based tester already deals with form errors of the mail address of the registration form.

With the same logic of adding actions, choices and data generation, we can cover all possibilities and combinations of our application. Using mutation testing we can also make sure that we do not overlook any straight forward possibility. The only coverage that will be left is errors originating from the database, the template engine and the implementation of the standard library. At first these seems out of reach but we can inject faults, e.g. database timeouts for statements, using mocks/fakes. However, there are more serious issues left to tackle. Even though we can now test all requirements and the complete source code of our application we are still not done: e.g. we do not have tested the security of our application.

# Bonus: Why Static Analysis does not always help?

In the previous section we talked about the problem that we now have tested our specification and using mutation testing we can even find all test cases that are missing from the perspective of our source code. However, since we are using third party packages we have to make sure that we are using them correctly, and also we need to make sure that we are using the programming language without introducing problems that cannot be found with mutation testing.

Both possibilities of problem sources are often presented as covered by static analysis tools. Unfortunately, this is not true as most static analysis tools are not strong enough to cover all possibilites and **proof** that code is free from a certain problem.

Since security is always an interesting issue (btw. we tried to put as many problems listed on https://owasp.org in the application as possible), we will show this problem with a great tool called [gosec](https://github.com/securego/gosec) which tries to find security issues by analyzing source code. However, as most static analysis tools gosec suffers from false positives. (**PLEASE NOTE:** we like gosec, it is just an example for the problems with static analysis tools) With the following command we run gosec on our application.

```bash
gosec ./comments/...
```

This command will output depending on the content of the repository something like the following.

```
[/go/src/github.com/symflower/sessions/2019/socrates-linz/comments/model/comment.go:29] - G202: SQL string concatenation (Confidence: HIGH, Severity: MEDIUM)
  > "INSERT INTO comments(mail, created, message) VALUES('" + mail


[/go/src/github.com/symflower/sessions/2019/socrates-linz/comments/model/user.go:23] - G202: SQL string concatenation (Confidence: HIGH, Severity: MEDIUM)
  > "SELECT mail, password FROM users WHERE mail = '"+mail


[/go/src/github.com/symflower/sessions/2019/socrates-linz/comments/model/user.go:35] - G202: SQL string concatenation (Confidence: HIGH, Severity: MEDIUM)
  > "INSERT INTO users(mail, password) VALUES('" + mail


[/go/src/github.com/symflower/sessions/2019/socrates-linz/comments/model/user.go:42] - G202: SQL string concatenation (Confidence: HIGH, Severity: MEDIUM)
  > "SELECT mail FROM users WHERE mail = '" + mail
```

Question: Are these real problems? Hint: If we look into the code, and maybe even test it, we can see that these are all real problems that can be formulated to SQL injections.

Finding such issues is great. However, the problem is that gosec (and other tools) report false positives. Let's introduce such a case. By changing the SQL statement in the file "comments/model/comment.go" of the function "CommentAll" to `table := "comments"; rows, err := db.Query("SELECT mail, created, message FROM " + table + " ORDER BY created DESC")`, we can run gosec again resulting into the following output.

```
[/go/src/github.com/symflower/sessions/2019/socrates-linz/comments/model/comment.go:37] - G202: SQL string concatenation (Confidence: HIGH, Severity: MEDIUM)
  > "SELECT mail, created, message FROM " + table
```

The reported warning is true: we are using a string concatentation for an SQL statement. However, in this case there is no problem with that at all.

Question: Why does the static analysis not know that this is not a problem? Hint: instead of defining a variable try `const table := "comments"` and run the tool again. Does this make a difference?

Hint: The problem is that this particular static analysis does not include the content of variables. Other static analysis include variables but give up on more complex code. Therefore, even good static analysis are not a perfect fit for finding such issues, since we always have to manually check if a warning is a true positive or a false positive.

Exercise: How can we add automatic SQL injection testing to our model-based tester? How can we add such testing without changing every form in our tester?

Exercise: Can we find other security issues in our application? (Hint: yes.) How can we find such problems automatically without adding the checks to every form, page, parameter, ...?

# Where does Symflower fit into all of this?

We https://symflower.com/ offer a product that analyses your source code and converts it into a mathematical model. This model is then used to calculate one specific unit test for every path including every problem and interesting combination. These test cases lead to a high code coverage, all without the need of any human input. Since we include specific problems in this analysis too, we automatically detect code quality issues and bugs such as overflows, memory corruptions and security issues. After such a test suite is generated, developers have only one step left to do: decide which test cases should be kept and which of the found problems should be fixed.

Sounds good? If you want to experience Symflower in action, checkout https://try.symflower.com.

# Problems in the application

- We have issues with parallel requests -> use the race detector.
- A lot can be refactored:
	- Use a common template for common data e.g. headers.
	- Model functions all have the fields of the structs as parameters. We can simply use an object of the struct as parameter.
	- We could use a struct as controller (or a context) which can then hold the database, request, response, ... as a state instead of adding new parameters to our interface.
	- We could use some form of middleware to return errors and return valid outputs instead of changing the response by hand.
- No IDs for User nor Post tables which makes addressing them more complicated.
- No primary/unique keys for User nor Post table which makes duplications possible.
- No foreign key between User and Post table.
- No locking for any table access which makes it possible to add a user twice under certain conditions.
- (?) The timestamp should be done with the database
- Errors should be not a magic constant.
- Validation is hardcoded in the controller code, should be moved to model.
- No mail length validation.
- No message length validation.
- Incorrect validation of email addresses.
	- There are validation libraries that can correctly validate email addresses. This is not as easy as it looks see https://tools.ietf.org/html/rfc3696
- The mail address is not validated in the comment form.
- No password validation.
- SQL injection of the mail, message and password.
- Password is not encrypted, not even hashed.
- XSS for mail address and message basically everywhere.
- Usage of static messages in the templates of the applications instead of using Internationalization  L18n features.
- Controller handles view (e.g. CSS and template).
- Forms all suffer of CSRF
- No encryption for the HTTP connection
- No specific timeouts for the HTTP server.
- Port is hardcoded
- No configuration for the HTTP servers port.
- Database is lost if the server dies.
- No configuration for the database.
- No cache for the user “login” for posting a message.
- No bruteforce limit for the “login” for posting a message.
- In general no quota for any access.

Did you find any other problems or can introduce new problems? Please tell us by creating a pull request to this repository.

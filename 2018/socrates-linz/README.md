# Advanced Testing Methods

This session gives you a short overview about some advanced testing methods. We start of with the basics of unit-testing of Go because this is then used for all other techniques. With mutation testing we can check that test cases really "test" something instead of just execute some code. With fuzzing we can really easily generate random test data or even whole test cases. With model-based testing we can create a model that can be used to generate targeted test data or even test cases. With delta debugging we can automatically reduce data, e.g. to find minimized data that still executes a certain error for a test case.

# Implementation

- Show "main.go" our implementation

# Unit Testing

- Unit testing is usually you get X in and check that Y comes out.
- Show "main_test.go our unit tests
- Run `go test -v -coverprofile=coverage.out` -> coverage: 58.3% of statements -> WHAT?!
- Run `go tool cover -html=coverage.out`
- Problem: You are just testing what you think is correct (e.g. only things that are defined in the specification).

# Mutation-testing

- https://en.wikipedia.org/wiki/Mutation_testing -> `Mutation testing (or mutation analysis or program mutation) is used to design new software tests and evaluate the quality of existing software tests. Mutation testing involves modifying a program in small ways. Each mutated version is called a mutant and tests detect and reject mutants by causing the behavior of the original version to differ from the mutant. This is called killing the mutant. Test suites are measured by the percentage of mutants that they kill.`
- https://github.com/zimmski/go-mutesting my implementation for Go
- Run `go-mutesting --debug github.com/symflower/sessions/2018/socrates-linz/csvler` -> `The mutation score is 0.333333 (3 passed, 6 failed, 3 duplicated, 0 skipped, total is 9)` -> There are just 3 "mutations" that are tested by our unit-tests but 6 cases are not! (3 are duplicates -> go-mutesting did some work that lead to the same program/coverage/results)

# Fuzzing

- https://en.wikipedia.org/wiki/Fuzzing -> `Fuzzing or fuzz testing is an automated software testing technique that involves providing invalid, unexpected, or random data as inputs to a computer program.`
- My implementation of a fuzzing/model-based testing and delta-debugging tool in one binary: https://github.com/zimmski/tavor The file format is defined here https://github.com/zimmski/tavor/blob/master/doc/format.md (there is also an API).
- Show "example.tavor"
- Run `tavor --format-file example.tavor fuzz`
- Run `tavor --format-file example.tavor fuzz --strategy AllPermutations`
- Show "csvler.tavor"
- Do the csvler.tavor with just two columns -> not likely to run into the bug
- Run `tavor --format-file csvler.tavor --max-repeat 10 fuzz` -> generates random CSV files with our format. And repeats things at a maximum 10 times.
- TODO There is a bug in Tavor where we cannot use +0,$cc.Value("," [1-9]) because the [1-9] item gets copied and not newly generated
- Run `bash fuzz.bash`

# Model-based testing

- https://en.wikipedia.org/wiki/Model-based_testing -> `Model-based testing is an application of model-based design for designing and optionally also executing artifacts to perform software testing or system testing. Models can be used to represent the desired behavior of a system under test (SUT), or to represent testing strategies and a test environment.`
- Model-based testing is not like fuzzing just about data... it is also about how you execute the system under test (SUT) with that data and how to validate that there is a problem or not.
- E.g. https://github.com/zimmski/tavor/blob/master/doc/complete-example.md -> coin(puts something in) and credit (checks that the credit is now correct)
- E.g. website checker: Home -> Login -> Profile -> Change Name (these are all actions on the website) the validation would then for instance check whether the response status is 200 OK.

# Delta Debugging

- https://en.wikipedia.org/wiki/Delta_debugging -> `Delta Debugging is a methodology to automate the debugging of programs using a scientific approach of hypothesis-trial-result loop.` -----> `In practice, the Delta Debugging algorithm builds on unit testing to isolate failure causes automatically `
- Show `delta-debugging-example.csv`
- Show `dd.tavor`
- TODO variables are a non-implemented feature in tavor for delta-debugging
- Run `tavor --verbose --format-file dd.tavor --max-repeat 20 reduce --input-file delta-debugging-example.csv --exec "go run main.go" --exec-exact-exit-code`

# The Future (Symflower)

- Show basic compare: https://symflower.com

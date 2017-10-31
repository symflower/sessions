package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/symflower/sessions/2017/socrates-linz-microservice-calculator"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var op socra.Operator
		if err := socra.Decode(r, &op); err != nil {
			panic(err)
		}

		out := calc(op)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		socra.Encode(w, socra.Result{out.Result})
	}).Methods("POST")

	socra.RunServer(r)
}

func calc(op socra.Operator) socra.Result {
	var a int
	var b int
	var err error

	if a, err = strconv.Atoi(op.A); err == nil {
		// Its a number!
	} else {
		o := parseOperator(op.A)
		a = calc(o).Result
	}

	if b, err = strconv.Atoi(op.B); err == nil {
		// Its a number!
	} else {
		o := parseOperator(op.B)
		b = calc(o).Result
	}

	switch op.Operator {
	case "+":
		return add(a, b)
	case "-":
		return sub(a, b)
	case "*":
		return mul(a, b)
	case "/":
		return div(a, b)
	default:
		panic("unkown op")
	}
}

func parseOperator(s string) socra.Operator {
	var op socra.Operator
	err := json.Unmarshal([]byte(s), &op)
	if err != nil {
		panic(err)
	}

	return op
}

func add(a int, b int) socra.Result {
	return socra.Post("http://add:8080/", socra.Add{A: a, B: b})
}

func sub(a int, b int) socra.Result {
	return socra.Post("http://sub:8080/", socra.Sub{A: a, B: b})
}

func mul(a int, b int) socra.Result {
	return socra.Post("http://mul:8080/", socra.Mul{A: a, B: b})
}

func div(a int, b int) socra.Result {
	return socra.Post("http://div:8080/", socra.Div{A: a, B: b})
}

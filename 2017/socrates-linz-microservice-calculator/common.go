package socra

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Operator struct {
	Operator string
	A        string
	B        string
}

type Add struct {
	A int
	B int
}

type Sub struct {
	A int
	B int
}

type Mul struct {
	A int
	B int
}

type Div struct {
	A int
	B int
}

type Result struct {
	Result int
}

func RunServer(r *mux.Router) {
	n := negroni.New()

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())

	n.UseHandler(r)

	n.Run(":8080")
}

func Decode(r *http.Request, v interface{}) error {
	d := json.NewDecoder(r.Body)

	err := d.Decode(v)
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	return nil
}

func Encode(w http.ResponseWriter, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func Post(url string, payload interface{}) Result {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(string(content))
	}

	var out Result
	if err := json.Unmarshal(content, &out); err != nil {
		panic(err)
	}

	return out
}

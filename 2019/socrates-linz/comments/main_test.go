package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestApplication(t *testing.T) {
	// Setup the server.
	go func() { // Run the server in the background.
		main()
	}()

	time.Sleep(time.Second) // Wait a second to let the server start.

	// Register a user.
	mail := "user@symflower.com"
	password := "secret"

	resp, err := http.PostForm("http://localhost:8080/register", url.Values{"mail": {mail}, "password": {password}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Log(string(body))
		t.Fatalf("Status code is %d", resp.StatusCode)
	}
	if !strings.Contains(string(body), "You registered the user "+mail) {
		t.Log(string(body))
		t.Fatal("Done message does not exist")
	}

	// Create a comment.
	message := "Some message!"

	resp, err = http.PostForm("http://localhost:8080/", url.Values{"mail": {mail}, "password": {password}, "message": {message}})
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Log(string(body))
		t.Fatalf("Status code is %d", resp.StatusCode)
	}
	if !strings.Contains(string(body), "We think comments are amazing!") {
		t.Log(string(body))
		t.Fatal("Application header does not exist")
	}
	if !strings.Contains(string(body), message) {
		t.Log(string(body))
		t.Fatal("Message does not exist")
	}
}

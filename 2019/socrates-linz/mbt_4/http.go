package mbt

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func httpHandleResponse(ctx *Context, resp *http.Response) (*http.Response, string) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Fatal(err)
	}
	resp.Body.Close()

	ctx.Body = string(body)

	ctx.FormData = map[string]string{}
	ctx.FormErrors = map[string]string{}

	return resp, ctx.Body
}

func HTTPGet(ctx *Context, uri string) (resp *http.Response, body string) {
	resp, err := http.Get("http://localhost:8080" + uri)
	if err != nil {
		ctx.Fatal(err)
	}

	return httpHandleResponse(ctx, resp)
}

func HTTPGetValid(ctx *Context, uri string) (resp *http.Response, body string) {
	resp, body = HTTPGet(ctx, uri)

	if resp.StatusCode != http.StatusOK {
		ctx.Fatalf("status code was %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	return resp, body
}

func HTTPPost(ctx *Context, uri string, data url.Values) (resp *http.Response, body string) {
	resp, err := http.PostForm("http://localhost:8080"+uri, data)
	if err != nil {
		ctx.Fatal(err)
	}

	return httpHandleResponse(ctx, resp)
}

func HTTPPostValid(ctx *Context, uri string, data url.Values) (resp *http.Response, body string) {
	resp, body = HTTPPost(ctx, uri, data)

	if resp.StatusCode != http.StatusOK {
		ctx.Fatalf("status code was %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	return resp, body
}

func HTTPPostDataSet(ctx *Context, key string, value string) {
	ctx.Writef("mbt.HTTPPostDataSet(ctx, %#v, %#v)\n", key, value)

	ctx.FormData[key] = value
}

func HTTPPostErrorSet(ctx *Context, key string, value string) {
	ctx.Writef("mbt.HTTPPostErrorSet(ctx, %#v, %#v)\n", key, value)

	ctx.FormErrors[key] = value
}

func HTTPPostSend(ctx *Context, uri string) (resp *http.Response, body string) {
	form := url.Values{}
	for k, v := range ctx.FormData {
		form[k] = []string{v}
	}

	resp, body = HTTPPostValid(ctx, uri, form)

	return resp, body
}

func DOM(ctx *Context, html string) *goquery.Document {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		ctx.Fatal(err)
	}

	return dom
}

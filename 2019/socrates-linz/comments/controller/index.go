package controller

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

var templateIndex = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Comments!</title>
		` + css + `
	</head>
	<body>
		<h1>We think comments are amazing!</h1>
		<form id="form_comment" class="pure-form pure-form-aligned" method="POST">
			<fieldset>
				<legend>Leave a comment or <a href="/register">register a user to comment</a></legend>

				<div class="pure-control-group">
					<label for="mail">Mail address</label>
					<input id="mail" name="mail" type="email" placeholder="Mail address" value="{{.Form.Mail}}">
					{{if .Form.MailError}}<span id="mail_error" class="pure-form-message-inline">{{.Form.MailError}}</span>{{end}}
				</div>

				<div class="pure-control-group">
					<label for="password">Password</label>
					<input id="password" name="password" type="password" placeholder="Password" value="{{.Form.Password}}">
					{{if .Form.PasswordError}}<span id="password_error" class="pure-form-message-inline">{{.Form.PasswordError}}</span>{{end}}
				</div>

				<div class="pure-control-group">
					<label for="message">Message</label>
					<input id="message" name="message" type="text" placeholder="Message" value="{{.Form.Message}}">
					{{if .Form.MessageError}}<span id="message_error" class="pure-form-message-inline">{{.Form.MessageError}}</span>{{end}}
				</div>

				<div class="pure-controls">
					<input type="submit" value="Create" class="pure-button pure-button-primary">
				</div>
			</fieldset>
		</form>
		<div class="comments">
			{{range $comment := .Comments}}
			<section>
				<div class="message">{{$comment.Message}}</div>
				<div class="by">By {{$comment.Mail}}</div>
			</section>
			{{end}}
		</div>
	</body>
</html>
`))

type templateIndexData struct {
	Comments []*model.Comment
	Form     struct {
		Mail          string
		MailError     string
		Password      string
		PasswordError string
		Message       string
		MessageError  string
	}
}

func HandleIndex(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	data := templateIndexData{}

	data.Form.Mail = r.FormValue("mail")
	data.Form.Password = r.FormValue("password")
	data.Form.Message = r.FormValue("message")

	if data.Form.Mail == "" {
		data.Form.MailError = "This is a required field."
	}

	if data.Form.Password == "" {
		data.Form.PasswordError = "This is a required field."
	} else {
		loggedin, err := model.UserLogin(db, data.Form.Mail, data.Form.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		} else if !loggedin {
			data.Form.PasswordError = "Incorrect login user or password"
		}
	}

	if data.Form.Message == "" {
		data.Form.MessageError = "This is a required field."
	}

	if data.Form.MailError == "" && data.Form.PasswordError == "" && data.Form.MessageError == "" {
		err := model.CommentAdd(db, data.Form.Mail, data.Form.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		data.Form.Mail = ""
		data.Form.Password = ""
		data.Form.Message = ""
	}

	var err error
	data.Comments, err = model.CommentAll(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = templateIndex.Execute(w, data)
	if err != nil {
		panic(err) // If we cannot render the template, the server is unusable.
	}
}

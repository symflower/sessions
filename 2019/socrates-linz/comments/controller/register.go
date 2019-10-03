package controller

import (
	"database/sql"
	"net/http"
	"strings"
	"text/template"

	"github.com/symflower/sessions/2019/socrates-linz/comments/model"
)

var templateRegister = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Comments!</title>
		` + css + `
	</head>
	<body>
		<h1>We think comments are amazing!</h1>
		{{if .Registered}}
			You registered the user {{.Form.Mail}}. You can now <a href="/">go commenting with your user.</a>
		{{else}}
			<form id="form_register" class="pure-form pure-form-aligned" method="POST">
				<fieldset>
					<legend>Register a user or <a href="/">go commenting with your user</a></legend>

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

					<div class="pure-controls">
						<input type="submit" value="Create" class="pure-button pure-button-primary">
					</div>
				</fieldset>
			</form>
		{{end}}
	</body>
</html>
`))

type templateRegisterData struct {
	Registered bool
	Form       struct {
		Mail          string
		MailError     string
		Password      string
		PasswordError string
	}
}

func HandleRegister(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	data := templateRegisterData{}

	data.Form.Mail = r.FormValue("mail")
	data.Form.Password = r.FormValue("password")

	if data.Form.Mail == "" {
		data.Form.MailError = "This is a required field."
	} else if !strings.Contains(data.Form.Mail, "@") {
		data.Form.MailError = "Not a valid mail address."
	} else if user, err := model.UserByMail(db, data.Form.Mail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	} else if user != nil {
		data.Form.MailError = "The mail address " + data.Form.Mail + " does already exist."
	}

	if data.Form.Password == "" {
		data.Form.PasswordError = "This is a required field."
	}

	if data.Form.MailError == "" && data.Form.PasswordError == "" {
		err := model.UserAdd(db, data.Form.Mail, data.Form.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		data.Registered = true
	}

	w.WriteHeader(http.StatusOK)
	err := templateRegister.Execute(w, data)
	if err != nil {
		panic(err) // If we cannot render the template, the server is unusable.
	}
}

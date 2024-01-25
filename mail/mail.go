package mail

import "html/template"

type Headers struct {
	From    string
	To      string
	Subject string
}

type Service interface {
	Ping() error
	Send(Headers, *template.Template, any) error
}

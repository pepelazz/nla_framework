package utils

import (
	"bytes"
	"github.com/go-gomail/gomail"
	"html/template"
)

func EmailSend(to, subject, emailBody string) error {

	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetAddressHeader("From", emailConfig.Sender, emailConfig.SenderName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(emailConfig.Host, int(emailConfig.Port), emailConfig.Sender, emailConfig.Password)

	return d.DialAndSend(m)
}

func EmailSendWithEmptySender(to, subject, emailBody string) error  {
	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetAddressHeader("From", emailConfig.Sender, emailConfig.SenderName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(emailConfig.Host, int(emailConfig.Port), "", emailConfig.Password)

	return d.DialAndSend(m)
}

func EmailSendChangePassword(to, href string) error  {
	data := struct {
		Name string
		Url string
	}{emailConfig.SenderName, href}

	t, err := template.New("letter").Parse(`
		<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
				"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html>
		</head>
		<body>
		<p>
			<h1>{{.Name}}</h1>
			<br>
			Для смены пароля кликните по ссылке<br>
			<a href="{{.Url}}">Смена пароля</a>
		</p>
		</body>
		</html>
`)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	emailBody := buf.String()

	if emailConfig.IsSendWithEmptySender {
		return EmailSendWithEmptySender(to, "Смена пароля", emailBody)
	}
	return EmailSend(to, "Смена пароля", emailBody)
}

func EmailSendRegistrationConfirm(to, href string) error  {
	data := struct {
		Name string
		Url string
	}{emailConfig.SenderName, href}

	t, err := template.New("letter").Parse(`
		<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
				"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html>
		</head>
		<body>
		<p>
			<h1>{{.Name}}</h1>
			<br>
			Для завершения процесса регистрации кликните по ссылке<br>
			<a href="{{.Url}}">Подтвердить регистрацию</a>
		</p>
		</body>
		</html>
`)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	emailBody := buf.String()

	if emailConfig.IsSendWithEmptySender {
		return EmailSendWithEmptySender(to, "Завершение процесса регистрации", emailBody)
	}

	return EmailSend(to, "Завершение процесса регистрации", emailBody)
}


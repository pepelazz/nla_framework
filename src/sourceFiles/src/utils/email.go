package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/matcornic/hermes"
	"github.com/pkg/errors"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

func EmailSendMessage(emails []string, subject, msg string, hermesEmail *hermes.Email) (err error) {

	if len(emailConfig.Sender) == 0 {
		errMsg := "Config.Email.Sender is nil. May be you forget write email sender in config.toml?"
		println(errMsg)
		return errors.New(errMsg)
	}

	var emailBody string
	emailBody, _, err = createEmailBody(subject, msg, hermesEmail)
	if err != nil {
		return
	}

	from := mail.Address{"", emailConfig.Sender}
	to := mail.Address{"", strings.Join(emails, ", ")}
	subj := subject
	//body := emailBody

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	//message := ""
	//for k, v := range headers {
	//	message += fmt.Sprintf("%s: %s\r\n", k, v)
	//}
	//message += "\r\n" + body

	// Connect to the SMTP Server
	servername := fmt.Sprintf("%s:%v", emailConfig.Host, emailConfig.Port)

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", emailConfig.Sender, emailConfig.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(emailBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}

func createEmailBody(subject, msg string, hermesEmail *hermes.Email) (emailBody string, MIME string, err error) {
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: emailConfig.SenderName,
			Link: webServerConfig.Url,
			// Optional product logo
			Logo:        emailConfig.SenderLogo,
			Copyright:   "",
			TroubleText: "Если возникли проблемы с действием'{ACTION}', скопируйте строку со ссылкой в браузер.",
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	if hermesEmail != nil {
		emailBody, err = h.GenerateHTML(*hermesEmail)
		if err != nil {
			fmt.Printf("GenerateHTML err %s\n", err)
			err = errors.Wrap(err, "GenerateHTML")
			return
		}
		MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	} else {
		//Generate the plaintext version of the e-mail (for clients that do not support xHTML)
		emailBody = msg
	}

	emailBody = fmt.Sprintf("From: %s <%s>\r\nSubject: %s\r\n%s\r\n%s", emailConfig.SenderName, emailConfig.Sender, subject, MIME, emailBody)
	return
}

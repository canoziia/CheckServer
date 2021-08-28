package mails

import (
	"strconv"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"
)

type MailConf struct {
	Port     string
	Host     string
	Encrypt  string
	Username string
	Password string
	Target   string
	Name     string
}

var (
	mailConf *MailConf
)

func init() {
	mailConf = new(MailConf)
}

func LoadMailConf(mconf MailConf) {
	*mailConf = mconf
}

// func SendMail(subject, msg, mailtype string) (err error) {
// 	auth := smtp.PlainAuth("", mailConf.Username, mailConf.Password, mailConf.Host)
// 	var content_type string
// 	if mailtype == "html" {
// 		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
// 	} else {
// 		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
// 	}

// 	body := []byte("To: " + mailConf.Target + "\r\nFrom: " + mailConf.Name + "<" + mailConf.Username + ">" + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + msg)
// 	targetList := strings.Split(mailConf.Target, ";")
// 	err = smtp.SendMail(mailConf.Host+":"+mailConf.Port, auth, mailConf.Username, targetList, body)
// 	return err
// }

func SendMail(subject, msg, mailtype string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", mailConf.Username)
	m.SetHeader("To", strings.Split(mailConf.Target, ";")...)
	// m.SetAddressHeader("Cc", "dan@dnlab.net", "Dan") 抄送
	m.SetHeader("Subject", subject)
	m.SetBody("text/"+mailtype, time.Now().Format("2006-01-02 15:04:05")+" "+msg)
	// m.Attach("/home/Alex/lolcat.jpg")
	port, err := strconv.Atoi(mailConf.Port)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(mailConf.Host, port, mailConf.Username, mailConf.Password)

	// Send the email to Bob, Cora and Dan.
	err = d.DialAndSend(m)
	return
}

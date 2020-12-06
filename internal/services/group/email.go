package group

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Solar-2020/Group-Backend/internal"
	"net/smtp"
	"sync"
	"text/template"
)

const (
	maxFails = 2
)

var (
	queue Queue
	queueInit sync.Once
)

type TemplateUser struct {
	FullName 	string
	Email 		string
	Avatar 		string
}

func sendInviteMessage(to string, adminName, adminSirname, adminEmail, adminAvatar string, group, link string) (err error) {
	from := internal.Config.InviteLetterSender
	host := internal.Config.InviteLetterHost

	tmplHTML, err := template.ParseFiles(internal.Config.InviteLetterBasePath + "/invite.html")
	if err != nil {
		return
	}
	tmplText, err := template.ParseFiles(internal.Config.InviteLetterBasePath + "/invite.txt")
	if err != nil {
		return
	}

	vars := struct{
		Admin TemplateUser
		GroupName string
		InviteLink string
	}{
		TemplateUser{
			FullName: adminName + " " + adminSirname,
			Email:    adminEmail,
			Avatar:   adminAvatar,
		},
		group,
		link,
	}

	htmlWriter := bytes.NewBufferString("")
	err = tmplHTML.Execute(htmlWriter, vars)
	if err != nil {
		return
	}

	textWriter := bytes.NewBufferString("")
	err = tmplText.Execute(textWriter, vars)
	if err != nil {
		return
	}

	varsBase := struct{
		To string
		From TemplateUser
		HtmlPart string
		TextPart string
	}{
		to,
		TemplateUser{
			FullName: "Pay Together",
			Email:    internal.Config.InviteLetterSender,
		},
		base64.StdEncoding.EncodeToString(htmlWriter.Bytes()),
		base64.StdEncoding.EncodeToString(textWriter.Bytes()),
	}

	tmplBase, err := template.ParseFiles(internal.Config.InviteLetterBasePath + "/invite.eml")
	if err != nil {
		return
	}
	baseWriter := bytes.NewBufferString("")
	err = tmplBase.Execute(baseWriter, varsBase)
	if err != nil {
		return
	}

	queueInit.Do(func(){
		queue.Init(internal.Config.InviteLetterTimespan)
	})

	queue.Enqueue(Package{
		Host:    host,
		From:    from,
		To:      to,
		Message: baseWriter.Bytes(),
	})

	go doSend()

	//auth := smtp.PlainAuth("", from, internal.Config.InviteLetterSenderPassword, host)
	//if err := smtp.SendMail(host+":25", auth, from, []string{to}, baseWriter.Bytes()); err != nil {
	//	return err
	//}

	return nil
}

func doSend() {
	p := queue.Dequeue()
	auth := smtp.PlainAuth("", p.From, internal.Config.InviteLetterSenderPassword, p.Host)
	if err := smtp.SendMail(p.Host+":25", auth, p.From, []string{p.To}, p.Message); err != nil {
		fmt.Println("Mailer: cannot send email to ", p.To, err)
		p.Retries +=1
		if p.Retries < maxFails {
			queue.Enqueue(p)
		}
		queue.timespan *= 3
		return
	}
	fmt.Println("Mailer: message sent to ", p.To)
	queue.timespan = internal.Config.InviteLetterTimespan
}
// code for u.s. export control compliance notifications
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/xoba/turd/taws"
)

/*

following the linux foundation's template at
https://www.linuxfoundation.org/export/

Project: OpenCue
Sent: 2019-05-18
    SUBMISSION TYPE: EAR 742.15(b) and 734.3(b)(3)
    SUBMITTED BY: Stephen Winslow
    SUBMITTED FOR: OpenCue Project a Series of LF Projects, LLC
    POINT OF CONTACT: Stephen Winslow
    TELEPHONE: 415-723-9709
    FAX: 415-723-9709
    MANUFACTURER(S): OpenCue Project a Series of LF Projects, LLC
    PRODUCT NAME/MODEL #: OpenCue
    ECCN: 5D002
    INTERNET LOCATION(S): https://github.com/imageworks/OpenCue; https://github.com/AcademySoftwareFoundation/

*/

type EarSubmission struct {
	Project      string
	Sent         time.Time
	Type         string
	By           string
	For          string
	Contact      string
	Phone        string
	Manufacturer string
	Product      string
	ECCN         string
	Locations    string
	Sender       string // name of person sending the email
}

type Config struct {
	DryRun  bool
	Prod    bool
	From    string
	Subject string
}

func main() {
	var c Config
	flag.BoolVar(&c.Prod, "p", false, "whether to run in production mode")
	flag.BoolVar(&c.DryRun, "d", true, "in dry run, don't send an email")
	flag.StringVar(&c.From, "from", "turd@xoba.com", "email sender")
	flag.StringVar(&c.Subject, "subject", "ear notification", "email subject")
	flag.Parse()
	if err := Run(c); err != nil {
		log.Fatal(err)
	}
}

func Run(c Config) error {
	const mike = "mike andrews"
	ear := EarSubmission{
		Project:      "turd",
		Sent:         time.Now(),
		Type:         "EAR 742.15(b) and 734.3(b)(3)",
		By:           mike,
		For:          "turd open source project",
		Contact:      mike,
		Phone:        "+19176086254",
		Manufacturer: mike + " and collaborators",
		Product:      "turd",
		ECCN:         "5D002",
		Locations:    "https://github.com/xoba/turd",
		Sender:       mike,
	}
	t, err := template.ParseFiles("docs/ear.template")
	if err != nil {
		return err
	}
	w := new(bytes.Buffer)
	if err := t.Execute(w, ear); err != nil {
		return err
	}
	var to string
	if c.Prod {
		to = "blah blah"
	} else {
		to = c.From
	}
	if c.DryRun {
		fmt.Println(w)
	} else {
		out, err := SendEmail(EmailParameters{
			From:    c.From,
			To:      strings.Split(to, ","),
			Subject: c.Subject,
			Body:    w.String(),
		})
		if err != nil {
			return err
		}
		fmt.Println(out)
		f, err := os.OpenFile("docs/receipts.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		e := json.NewEncoder(f)
		e.SetEscapeHTML(false)
		e.SetIndent("", "  ")
		if err := e.Encode(out); err != nil {
			return err
		}
	}
	return nil
}

type EmailParameters struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type Receipt struct {
	EmailParameters
	MessageID string
}

func (r Receipt) String() string {
	buf, _ := json.MarshalIndent(r, "", "  ")
	return string(buf)
}

func SendEmail(e EmailParameters) (*Receipt, error) {
	session, err := taws.NewSession()
	if err != nil {
		return nil, err
	}
	content := func(s string) *ses.Content {
		return &ses.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(s),
		}
	}
	addrs := func(list []string) (out []*string) {
		for _, x := range list {
			x = strings.TrimSpace(x)
			x = strings.ToLower(x)
			if len(x) == 0 {
				continue
			}
			out = append(out, aws.String(x))
		}
		return
	}
	resp, err := ses.New(session).SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  addrs(e.To),
			BccAddresses: []*string{},
		},
		Message: &ses.Message{
			Subject: content(e.Subject),
			Body: &ses.Body{
				Html: content(e.Body),
			},
		},
		Source: aws.String(e.From),
	})
	if err != nil {
		return nil, err
	}
	r := Receipt{
		EmailParameters: e,
		MessageID:       *resp.MessageId,
	}
	return &r, nil
}

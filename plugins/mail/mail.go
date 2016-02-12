package mail

import (
	"strconv"

	"./../modules"
	"./../settings"
	"github.com/robertkrimen/otto"
	"gopkg.in/gomail.v1"
)

type SMTPSettings struct {
	Username   string
	Password   string
	Servername string
	Port       int
}

var config SMTPSettings

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "mail",
		Init: func(vm *otto.Otto) otto.Value {
			o, _ := vm.Object("({})")

			o.Set("loadMailSettings", func(c otto.FunctionCall) otto.Value {
				loadSettings(vm)
				return otto.TrueValue()
			})
			o.Set("send", func(c otto.FunctionCall) otto.Value {
				loadSettings(vm)
				recv, _ := c.Argument(0).ToString()
				subject, _ := c.Argument(1).ToString()
				msg, _ := c.Argument(2).ToString()

				err := sendmail(recv, subject, msg, "")

				return modules.ToResult(vm, true, err)
			})
			return o.Value()
		},
	}

	return &p1
}

func loadSettings(vm *otto.Otto) {
	settings := settings.GetSettings()

	mailSettings := settings.Plugins["mail"]

	uname := mailSettings["username"]
	pw := mailSettings["password"]
	sname := mailSettings["servername"]
	portStr := mailSettings["port"]

	config = SMTPSettings{
		Username:   uname,
		Password:   pw,
		Servername: sname,
	}

	port, _ := strconv.Atoi(portStr)
	config.Port = port
}

func sendmail(email string, subject string, messageString string, txtAttachment string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Username)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", messageString)
	if len(txtAttachment) > 0 {
		f := gomail.CreateFile("attached.txt", []byte(txtAttachment))
		msg.Attach(f)
	}

	mailer := gomail.NewMailer(config.Servername, config.Username, config.Password, config.Port)
	err := mailer.Send(msg)
	return err
}

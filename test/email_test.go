package test

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "from <949244762@qq.com>"  // 发送者
	e.To = []string{"949244762@qq.com"} // 接收者
	e.Subject = "test email 验证"
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "949244762@qq.com", "qtssdidxkuytbcah", "smtp.qq.com"))
	//mtkubqxqhwfxbbff
	//err := e.SendWithTLS("smtp.qq.com:587", smtp.PlainAuth("", "Aurora@gmail.com", "qtssdidxkuytbcah", "smtp.qq.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}

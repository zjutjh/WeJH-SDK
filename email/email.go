package email

import (
	"gopkg.in/gomail.v2"
)

// MailboxConf 发送邮件内容配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
}

// InfoConfig 发件人配置
type InfoConfig struct {
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

// SendEmail 一个message可以是这样的：
// message := fmt.Sprintf(`<div>
//
//	    <div>
//	        hello
//	    </div>
//	</div>`)
func (config InfoConfig) SendEmail(message string, mailbox MailboxConf) error {
	m := gomail.NewMessage()
	m.SetHeader(`From`, config.Sender)
	m.SetHeader(`To`, mailbox.RecipientList...)
	m.SetHeader(`Subject`, mailbox.Title)
	m.SetBody(`text/html`, message)
	// m.Attach("./Dockerfile") //添加附件
	d := gomail.NewDialer(config.SMTPAddr, config.SMTPPort, config.Sender, config.SPassword)
	return d.DialAndSend(m)
}

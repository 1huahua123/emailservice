package email

import (
	"bufio"
	"bytes"
	"emailservice/config"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

// EmailService 结构体表示电子邮件服务
type EmailService struct {
	config *config.Config
}

// NewEmailService 是 EmailService 的构造函数
func NewEmailService(config *config.Config) *EmailService {
	return &EmailService{config: config}
}

// Email 结构体表示电子邮件
type Email struct {
	To          []string // 收件人
	Cc          []string // 抄送人
	Bcc         []string // 密送人
	Subject     string   // 主题
	Body        string   // 正文
	Attachments []string // 附件路径
}

// SendEmail 方法用于发送电子邮件
func (s *EmailService) SendEmail(email *Email) error {
	from := s.config.SMTP.Username // 发件人
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	// 构建邮件消息
	msg, err := s.buildMessage(from, email)
	// SMTP 地址
	smtpAddress := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	// 发送邮件
	err = smtp.SendMail(smtpAddress, auth, from, append(email.To, append(email.Cc, email.Bcc...)...), msg)
	if err != nil {
		return err
	}
	return nil
}

// buildMessage 方法构建邮件消息
func (s *EmailService) buildMessage(from string, email *Email) ([]byte, error) {
	var buffer bytes.Buffer
	boundary := "my-boundary-123" // 分隔符，用于区分邮件的不同部分
	header := make(map[string]string)
	header["From"] = from
	header["To"] = strings.Join(email.To, ",")
	if len(email.Cc) > 0 {
		header["Cc"] = strings.Join(email.Cc, ",")
	}
	header["Subject"] = email.Subject
	header["MIME-Version"] = "1.0"
	header["Content-type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)
	// 写入头信息
	for k, v := range header {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buffer.WriteString("\r\n")
	// 定义一个写入分隔符的函数
	writeBoundary := func() {
		buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}

	// 写入邮件正文
	writeBoundary()
	buffer.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n\r\n")
	buffer.WriteString(email.Body)
	buffer.WriteString("\r\n")

	// 写入附件
	for _, attachment := range email.Attachments {
		err := s.writeAttachment(&buffer, attachment, boundary)
		if err != nil {
			return nil, err
		}
	}

	// 结束分隔符
	buffer.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return buffer.Bytes(), nil
}

// writeAttachment 方法写入附件内容
func (s *EmailService) writeAttachment(buffer *bytes.Buffer, attachment, boundary string) error {
	file, err := os.Open(attachment)
	if err != nil {
		return fmt.Errorf("打开附件时出错 %s: %w", attachment, err)
	}
	defer file.Close()

	writeBoundary := func() {
		buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}

	writeBoundary()
	buffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filepath.Base(attachment)))
	buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
	buffer.WriteString("Content-Type: application/octet-stream\r\n\r\n")

	writer := base64.NewEncoder(base64.StdEncoding, buffer)
	defer writer.Close()

	reader := bufio.NewReader(file)
	buf := make([]byte, 3*1024) // 定义读取缓冲区
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取附件时出错 %s: %w", attachment, err)
		}
		if n == 0 {
			break
		}
		if _, err := writer.Write(buf[:n]); err != nil {
			return fmt.Errorf("编码附件时出错 %s: %w", attachment, err)
		}
	}
	buffer.WriteString("\r\n")

	return nil
}

// ParseTemplate 方法解析模板文件并返回渲染后的字符串
func (s *EmailService) ParseTemplate(templateFileName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

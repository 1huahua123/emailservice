package main

import (
	"emailservice/config"
	"emailservice/email"
	"log"
	"time"
)

func main() {
	// 加载配置文件
	config.LoadConfig()
	// 创建一个新的电子邮件服务实例
	emailService := email.NewEmailService(config.AppConfig)
	// 电子邮件模板数据
	emailData := map[string]interface{}{
		"Title":  "Test Email",            // 标题
		"Header": "Welcome!",              // 头部
		"Body":   "This is a test email.", // 主体内容
	}

	// 解析模板文件
	var err error
	body, err := emailService.ParseTemplate("templates/example.html", emailData)
	if err != nil {
		log.Fatalf("分析电子邮件模板时出错: %s", err)
	}
	// 创建一个新的电子邮件实例
	email := &email.Email{
		To:          []string{"1850529744@qq.com"},      // 收件人
		Cc:          []string{"1850529744@qq.com"},      // 抄送人
		Bcc:         []string{"1850529744@qq.com"},      // 密送人
		Subject:     "Test Email",                       // 主题
		Body:        body,                               // 正文
		Attachments: []string{"path/to/attachment.txt"}, // 附件路径
	}
	// 重试次数
	retries := 3
	// 重试发送电子邮件
	for i := 0; i < retries; i++ {
		err = emailService.SendEmail(email)
		if err == nil {
			log.Println("电子邮件发送成功！")
			break
		}

		log.Printf("发送电子邮件失败 (尝试重发 %d/%d): %v", i+1, retries, err)
		time.Sleep(2 * time.Second) // 等待2s再重试
	}
	// 如果多次尝试后仍然失败，记录错误日志并终止程序
	if err != nil {
		log.Fatalf("尝试 %d 次后发送电子邮件失败: %v", retries, err)
	}
}

# 电子邮件服务

这是一个用 Go 编写的电子邮件服务，支持通过配置文件加载 SMTP 配置，并可以发送带有自定义模板和附件的电子邮件。该服务具有良好的扩展性，易于维护和扩展。

## 功能
- 支持通过配置文件加载 SMTP 配置。
- 支持发送带有 HTML 模板的电子邮件。
- 支持发送带有附件的电子邮件。
- 提供重试机制，发送失败时会进行多次重试。

## 目录结构
```
├── config
│ ├── config.go
├── email
│ ├── email.go
├── templates
│ └── example.html
├── main.go
├── go.mod
├── config.yaml
└── README.md
```

## 配置文件

在 `config/config.yaml` 中配置 SMTP 服务器的相关信息：

```yaml
smtp:
  host: "smtp.example.com"
  port: 587
  username: "your_username"
  password: "your_password"
```

## 模版文件

在 `templates` 目录中创建 HTML 模板文件，例如 `example.html`：

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Header}}</h1>
    <p>{{.Body}}</p>
</body>
</html>
```

## 使用方法

### 1.安装依赖
在你的项目目录中初始化 Go 模块，并安装所需依赖：
```shell
go mod init yourproject
go get github.com/spf13/viper
```
### 2.导入包
```go
import (
    "emailservice/config"
    "emailservice/email"
)
```

### 3.加载配置
首先加载配置文件：
```go
config.LoadConfig()
```
### 4.创建 EmailService 实例
使用加载的配置创建一个 EmailService 实例：
```go
emailService := email.NewEmailService(config.AppConfig)
```

### 5.解析模板
使用模板文件和数据创建邮件正文：
```go
emailData := map[string]interface{}{
    "Title":  "Test Email",
    "Header": "Welcome!",
    "Body":   "This is a test email.",
}
body, err := emailService.ParseTemplate("templates/example.html", emailData)
if err != nil {
    log.Fatalf("分析电子邮件模板时出错: %s", err)
}
```

### 6.创建 Email 实例
定义电子邮件的收件人、主题、正文和附件路径：
```go
email := &email.Email{
    To:          []string{"recipient@example.com"},
    Subject:     "Test Email",
    Body:        body,
    Attachments: []string{"path/to/attachment.txt"},
}
```

### 7.发送电子邮件
使用重试机制发送电子邮件：
```go
retries := 3

for i := 0; i < retries; i++ {
    err = emailService.SendEmail(email)
    if err == nil {
        log.Println("电子邮件发送成功！")
        break
    }

    log.Printf("发送电子邮件失败 (尝试重发 %d/%d): %v", i+1, retries, err)
    time.Sleep(2 * time.Second) // 等待一段时间再重试
}

if err != nil {
    log.Fatalf("尝试 %d 次后发送电子邮件失败: %v", retries, err)
}
```

## 完整示例
```go
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
        "Title":  "Test Email",
        "Header": "Welcome!",
        "Body":   "This is a test email.",
    }

    // 解析模板文件
    var err error
    body, err := emailService.ParseTemplate("templates/example.html", emailData)
    if err != nil {
        log.Fatalf("分析电子邮件模板时出错: %s", err)
    }

    // 创建一个新的电子邮件实例
    email := &email.Email{
        To:          []string{"recipient@example.com"},
        Subject:     "Test Email",
        Body:        body,
        Attachments: []string{"path/to/attachment.txt"},
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
        time.Sleep(2 * time.Second) // 等待一段时间再重试
    }

    // 如果多次尝试后仍然失败，记录错误日志并终止程序
    if err != nil {
        log.Fatalf("尝试 %d 次后发送电子邮件失败: %v", retries, err)
    }
}
```

## 依赖
- Viper：用于加载和管理配置文件。
- Go 标准库的 net/smtp：用于发送电子邮件。

### 如果你有任何问题或建议，请随时提出。谢谢！
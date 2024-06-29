package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config 结构体表示配置文件的顶层结构
type Config struct {
	SMTP SMTPConfig `mapstructure:"smtp"` // SMTP 字段包含邮件服务器的配置
}

// SMTPConfig 结构体表示 SMTP 服务器的配置项
type SMTPConfig struct {
	Host     string `mapstructure:"host"`     // 服务器地址
	Port     int    `mapstructure:"port"`     // 服务器端口
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
}

// AppConfig 是一个全局变量，用于存储应用程序的配置
var AppConfig *Config

// LoadConfig 函数用于加载和解析配置文件
func LoadConfig() {
	viper.SetConfigName("config") // 设置配置文件的名称（不包括扩展名）
	viper.SetConfigType("yaml")   // 设置配置文件的类型
	viper.AddConfigPath(".")      // 设置配置文件的路径，这里表示当前目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件时出错: %s", err)
	}
	// 初始化 AppConfig
	AppConfig = &Config{}
	// 将配置文件中的内容反序列化到 AppConfig 结构体中
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("解组配置时出错: %s", err)
	}
}

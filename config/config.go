package config

import (
	"gopkg.in/ini.v1"
	"log"
	"time"
)

// App App配置结构体
type App struct {
	RunMode   string
	JwtSecret string
}

// AppConf app配置
var AppConf = &App{}

// Server 服务器配置映射
type Server struct {
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ServerConf 服务器配置
var ServerConf = &Server{}

// Database 数据库配置映射
type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

// DatabaseConf 数据库配置
var DatabaseConf = &Database{}

// 设置配置 path: 配置文件的路径
func Setup(path string) {
	Cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("Failed to parse app.ini:%v", err)
	}
	err = Cfg.Section("app").MapTo(AppConf)
	if err != nil {
		log.Fatalf("Section 'app' mapping err: %v", err)
	}
	err = Cfg.Section("server").MapTo(ServerConf)
	if err != nil {
		log.Fatalf("Section 'server' mapping err: %v", err)
	}

	ServerConf.ReadTimeout = ServerConf.ReadTimeout * time.Second
	ServerConf.WriteTimeout = ServerConf.WriteTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseConf)
	if err != nil {
		log.Fatalf("Section 'database' mapping err: %v", err)
	}
}

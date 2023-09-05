package config

import (
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"time"
)

type AppConf struct {
	AppName  string    `yaml:"app_name"`
	Server   *Server   `yaml:"server"`
	Logger   *Logger   `yaml:"logger"`
	DB       *DB       `yaml:"database_url"`
	Token    *Token    `yaml:"token"`
	Telegram *Telegram `yaml:"telegram"`
	ChatGPT  *ChatGPT  `yaml:"chatgpt"`
}

type ChatGPT struct {
	Token string `yaml:"token"`
}

type Telegram struct {
	AppID   int    `yaml:"app_id"`
	ApiHash string `yaml:"token"`
	Phone   string `yaml:"phone"`
}

type DB struct {
	Name     string `yaml:"name"`
	User     string `json:"-" yaml:"user"`
	Password string `json:"-" yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Token struct {
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Server struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

func NewAppConf() AppConf {
	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		Logger: &Logger{
			Level: os.Getenv("LOG_LEVEL"),
		},
		Server: &Server{
			Port: os.Getenv("SERVER_PORT"),
		},
		DB: &DB{
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Telegram: &Telegram{
			AppID:   getenvInt("TELEGRAM_APP_ID"),
			ApiHash: os.Getenv("TELEGRAM_API_HASH"),
			Phone:   os.Getenv("TELEGRAM_PHONE"),
		},
		ChatGPT: &ChatGPT{
			Token: os.Getenv("CHAT_GPT_TOKEN"),
		},
		Token: &Token{
			AccessTTL:     time.Duration(getenvInt("ACCESS_TTL")) * time.Minute,
			RefreshTTL:    time.Duration(getenvInt("REFRESH_TTL")) * time.Minute,
			AccessSecret:  os.Getenv("ACCESS_SECRET"),
			RefreshSecret: os.Getenv("REFRESH_SECRET"),
		},
	}
}

func getenvInt(key string) int {
	env := os.Getenv(key)

	envInt, err := strconv.Atoi(env)
	if err != nil {
		log.Panicln(err)
	}
	return envInt

}

func (a *AppConf) Init(logger *zap.Logger) {
	shutDownTimeOut, err := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse server shutdown timeout error")
	}
	shutDownTimeout := time.Duration(shutDownTimeOut) * time.Second
	if err != nil {
		logger.Fatal("config: parse rpc server shutdown timeout error")
	}

	a.Server.ShutdownTimeout = shutDownTimeout
}

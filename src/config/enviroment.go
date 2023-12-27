package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        int    `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPass        string `mapstructure:"DB_PASS"`
	DBName        string `mapstructure:"DB_NAME"`
	Email         string `mapstructure:"EMAIL"`
	Password      string `mapstructure:"PASSWORD"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      int    `mapstructure:"SMTP_PORT"`
	EmailBody     string `mapstructure:"EMAIL_BODY"`
	Subject       string `mapstructure:"SUBJECT"`
	Attachments   string `mapstructure:"ATACHMENTS"`
	Begin         int    `mapstructure:"BEGIN"`
	EmailsForPack int    `mapstructure:"EMAILS_FOR_PACK"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}

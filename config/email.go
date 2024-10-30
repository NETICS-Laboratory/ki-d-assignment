package config

import "ki-d-assignment/utils"

type EmailConfig struct {
	Host         string `mapstructure:"SMTP_HOST"`
	Port         int    `mapstructure:"SMTP_PORT"`
	SenderName   string `mapstructure:"SMTP_SENDER_NAME"`
	AuthEmail    string `mapstructure:"SMTP_AUTH_EMAIL"`
	AuthPassword string `mapstructure:"SMTP_AUTH_PASSWORD"`
}

func NewEmailConfig() (*EmailConfig, error) {
	var config EmailConfig
	config.Host = utils.MustGetenv("SMTP_HOST")
	config.Port = utils.MustGetenvInt("SMTP_PORT")
	config.SenderName = utils.MustGetenv("SMTP_SENDER_NAME")
	config.AuthEmail = utils.MustGetenv("SMTP_AUTH_EMAIL")
	config.AuthPassword = utils.MustGetenv("SMTP_AUTH_PASSWORD")

	return &config, nil
}

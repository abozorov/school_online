package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP             http
		UserService      userService
		ClassroomService classroomService
		RaitingService   raitingService
		PG               pg
		JWT              jwt
		Email            email
	}

	http struct {
		// server
		Port            string `env:"SERVER_PORT"`
		Storage         string `env:"STORAGE"`
		AuditLogStorage string `env:"AUDIT_LOG_STORAGE"`
	}

	userService struct {
		Port string `env:"USER_SERVICE_PORT"`
		Host string `env:"USER_SERVICE_HOST"`
	}

	classroomService struct {
		Port string `env:"CLASSROOM_SERVICE_PORT"`
		Host string `env:"CLASSROOM_SERVICE_HOST"`
	}

	raitingService struct {
		Port string `env:"RAITING_SERVICE_PORT"`
		Host string `env:"RAITING_SERVICE_HOST"`
	}

	pg struct {
		// database
		Host     string `env:"PG_HOST"`
		Port     string `env:"PG_PORT"`
		User     string `env:"PG_USER"`
		Password string `env:"PG_PASSWORD"`
		Name     string `env:"PG_NAME"`
	}

	jwt struct {
		SecretToken string `env:"JWT_SECRET"`
		JWTLiveTime int    `env:"JWT_LIVE_TIME"`
	}

	email struct {
		Email    string `env:"EMAIL"`
		Password string `env:"EMAIL_PASSWORD"`
		Host     string `env:"EMAIL_HOST"`
		Port     string `env:"EMAIL_PORT"`
	}
)

func NewConfig(path string) (*Config, error) {
	var cnf Config

	err := cleanenv.ReadConfig(path, &cnf)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return &cnf, nil
}

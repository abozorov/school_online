package userservice

import (
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	"github.com/abozorov/school_online/pkg/cache"
	"github.com/abozorov/school_online/pkg/jwt"
	mailsender "github.com/abozorov/school_online/pkg/mail_sender"
)

type UserService struct {
	serviceManager services.IServiceManager
	jwt            *jwt.JWTSecret
	memCache       cache.ICache
	mailSender     *mailsender.MailSender
}

func NewUserService(
	serviceManager services.IServiceManager,
	jwt *jwt.JWTSecret,
	memCache cache.ICache,
	mailsender *mailsender.MailSender) *UserService {

	return &UserService{
		serviceManager: serviceManager,
		jwt:            jwt,
		memCache:       memCache,
		mailSender:     mailsender,
	}
}

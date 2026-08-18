package middleware

import "github.com/abozorov/school_online/pkg/jwt"

type Middleware struct {
	jwt *jwt.JWTSecret
}

func NewMiddlware(jwt *jwt.JWTSecret) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}

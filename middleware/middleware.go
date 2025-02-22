package middleware

import (
	"log"
	"time"

	"github.com/PhantomX7/dhamma/config"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	authMiddleware *jwt.GinJWTMiddleware
	// enforcer       *casbin.Enforcer
}

func New() *Middleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte(config.JWT_SECRET),
		Timeout:    24 * time.Hour,
		MaxRefresh: 6 * time.Hour,
		TimeFunc:   time.Now,
	})
	if err != nil {
		log.Fatal("jwt-error:" + err.Error())
	}

	return &Middleware{
		authMiddleware: authMiddleware,
	}
}

func (m *Middleware) AuthHandle() gin.HandlerFunc {
	return m.authMiddleware.MiddlewareFunc()
}

func (m *Middleware) RefreshHandle() gin.HandlerFunc {
	return m.authMiddleware.RefreshHandler
}

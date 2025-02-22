package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
)

type Middleware struct {
	authMiddleware *jwt.GinJWTMiddleware
	enforcer       *casbin.Enforcer
}

package middleware

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
	"gorm.io/gorm"

	"github.com/PhantomX7/go-core/utility/errors"
)

type Config struct {
	// put middleware config here
	JwtToken string
	Db       *gorm.DB
	Enforcer *casbin.Enforcer // casbin enforcer
}

type Middleware struct {
	config         Config
	authMiddleware *jwt.GinJWTMiddleware
	enforcer       *casbin.Enforcer
}

func New(cfg Config) *Middleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:              []byte(cfg.JwtToken),
		Timeout:          6 * time.Hour,
		MaxRefresh:       6 * time.Hour,
		TimeFunc:         time.Now,
		SigningAlgorithm: "HS512",
	})
	if err != nil {
		log.Fatal("jwt-error:" + err.Error())
	}

	return &Middleware{
		config:         cfg,
		authMiddleware: authMiddleware,
		enforcer:       cfg.Enforcer,
	}
}

func (m *Middleware) AuthHandle() gin.HandlerFunc {
	return m.authMiddleware.MiddlewareFunc()
}

func (m *Middleware) RefreshHandle() gin.HandlerFunc {
	return m.authMiddleware.RefreshHandler
}

func (m *Middleware) LanguageHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		language := c.GetHeader("Language")
		if language == "" {
			c.Request.Header.Add("Language", "id")
		}
		c.Next()
	}
}

func (m *Middleware) ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Find out what type of error it is

				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						// check if it is part of custom error
						if err, ok := e.Err.(errors.CustomError); ok {
							// print the underlying error and return the specified message to user
							c.JSON(err.HTTPCode, gin.H{
								"errors":  nil,
								"message": err.Message,
							})
						} else {
							c.JSON(c.Writer.Status(), gin.H{
								"errors":  nil,
								"message": e.Error(),
							})
						}

					}
				case gin.ErrorTypeBind:
					errs, ok := e.Err.(validator.ValidationErrors)
					if ok {
						list := make(map[string]string)
						for _, err := range errs {
							list[strcase.SnakeCase(err.Field())] = validationErrorToText(err)
						}

						// Make sure we maintain the preset response status
						status := http.StatusUnprocessableEntity
						if c.Writer.Status() != http.StatusOK {
							status = c.Writer.Status()
						}
						c.JSON(status, gin.H{
							"errors":  list,
							"message": "validation error",
						})
					} else {
						c.JSON(422, gin.H{
							"errors":  nil,
							"message": "please make sure to provide the correct data type or format",
						})
					}
				default:
					// Log all other errors
					//rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
				}

			}
			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"errors": "something went wrong"})
			}
		}
	}
}

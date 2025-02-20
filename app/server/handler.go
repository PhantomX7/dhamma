package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/PhantomX7/go-core/utility/validators"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/PhantomX7/dhamma/app/middleware"
)

// Handler All handler that need to be registered MUST implement this interface
type Handler interface {
	Register(r *gin.Engine, m *middleware.Middleware)
}

// BuildHandler will build http handler with giver middleware and all handlers
func BuildHandler(middleware *middleware.Middleware, handlers ...Handler) http.Handler {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// register all custom validator here
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("unique",
			validators.CustomValidator.Unique())
		if err != nil {
			log.Println("error when applying unique validator")
		}
		err = v.RegisterValidation("exist", validators.CustomValidator.Exist())
		if err != nil {
			log.Println("error when applying exist validator")
		}
	}

	// apply global middleware here
	config := cors.Config{
		//AllowOrigins: strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		//AllowWebSockets:  true,
		//AllowCredentials: true,
		AllowAllOrigins: true,
		AllowWildcard:   true,
		AllowHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Bearer",
			"Connection",
			"Cookie",
			"Content-Length",
			"Content-Type",
			"Origin",
			"Authorization",
			"X-Forwarded-For",
			"X-Real-Ip",
			"User-Agent",
			"Lang",
			"Version",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Access-Control-Allow-Origin",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	}

	router.Use(cors.New(config))
	router.Use(middleware.ErrorHandle())

	// set max upload file size
	//router.MaxMultipartMemory = 8 << 20  // 8 MiB

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/healthz", healthz)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//router.GET("/test",test)

	// start registering routes from all handlers
	for _, reg := range handlers {
		reg.Register(router, middleware)
	}

	// 404 not found function
	router.NoRoute(notFound)

	return router
}

// healthz godoc
// @Summary      Check server health status
// @Description  will return ok with number of go routine running
// @Tags         handler
// @Accept       json
// @Produce      json
// @Success      200  {string}   string  "ok"
// @Failure      400  {string}   string  "ok"
// @Failure      404  {string}   string  "ok"
// @Failure      500  {string}   string  "ok"
// @Router       /healhtz [get]
func healthz(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprint("ok:", runtime.NumGoroutine()))
}

func notFound(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotFound)
}

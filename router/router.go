package router

import (
	"github.com/parikhrahil/go-micro-starter/dto"
	"github.com/parikhrahil/go-micro-starter/logger"
	"github.com/parikhrahil/go-micro-starter/middleware"

	"github.com/gin-gonic/gin"
)

type Opts struct {
	Logger logger.Logger
	Mode   string
}

func New(opts *Opts) *gin.Engine {
	setMode(opts.Mode)
	router := gin.New()

	// Add middleware for logging and error handling
	logging := middleware.LogHandler(opts.Logger)
	error := middleware.ErrorHandler()

	router.Use(logging)
	router.Use(error)

	// Register ping route
	registerPing(router)
	return router
}

func setMode(mode string) {
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func registerPing(router *gin.Engine) {
	router.GET("/ping", func(ctx *gin.Context) {
		data := map[string]string{"message": "pong"}
		dto.OK(ctx, data)
	})
}

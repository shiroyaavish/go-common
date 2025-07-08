package http_server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strings"

	"github.com/IntelXLabs-LLC/go-common/config"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

type Router struct {
	*fiber.App
	hostname   string
	isProd     bool
	port       int
	srvStarted bool
}

type RouterConfig struct {
	AppName           string         `json:"app_name"`
	StrictRouting     bool           `json:"strict_routing"`
	IsRecoveryEnabled bool           `json:"is_recovery_enabled"`
	IsLoggerEnabled   bool           `json:"is_logger_enabled"`
	CompressionType   compress.Level `json:"compression_type"`
	Port              int            `json:"port"`
	GlobalCors        []string       `json:"global_cors"`
}

func NewRouter(cfg RouterConfig) *Router {
	fiberConfig := fiber.Config{}

	fiberConfig.AppName = cfg.AppName
	fiberConfig.StrictRouting = cfg.StrictRouting

	h := &Router{
		App:  fiber.New(fiberConfig),
		port: cfg.Port,
	}

	if cfg.IsRecoveryEnabled {
		h.Use(fiberRecover.New(fiberRecover.Config{
			EnableStackTrace: true,
		}))
	}

	if cfg.IsLoggerEnabled {
		h.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
			// Don't log health check. (To avoid ALB Health check Spam)
			return strings.HasPrefix(c.Path(), "api")
		}}))
	}

	corsHosts := ""
	if config.GetCurrentEnvironment() == config.Production {
		h.isProd = true
	} else {
		h.isProd = false
	}

	h.Use(cors.New(cors.Config{
		// Allowed Origins are based on environment.
		AllowOrigins: corsHosts,
		AllowHeaders: "*",
	}))

	h.Use(compress.New(compress.Config{
		Level: cfg.CompressionType,
	}))

	if hostname, _ := os.Hostname(); len(hostname) > 0 {
		h.hostname = hostname
	}

	return h
}

func (r *Router) StartAsync() {
	if r.srvStarted {
		return
	}

	log.Println("Starting server on port", r.port)
	r.srvStarted = true
	go func() {
		err := r.Listen(fmt.Sprintf(":%d", r.port))
		if err != nil {
			panic(err)
		}
	}()
}

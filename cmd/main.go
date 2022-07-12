package main

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"github.com/umi0410/umi-mrdebugger/handler"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Post("/", handler.HelpDebugging)
	app.Get("/health", handler.Health)
	log.Info("server is listening to 0.0.0.0:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Error(err)
	}
	//http.HandleFunc("/currency", handler.ServeHTTP)
	//http.HandleFunc("/health", handler.HealthCheck)

	//http.ListenAndServe(":8080", nil)
}

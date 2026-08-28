package main

import (
	"log"

	"github.com/akhtarfath/config"
	"github.com/akhtarfath/routes"
	"github.com/gofiber/fiber/v3"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	app := fiber.New()

	routes.New(app)

	addr := ":" + config.Port()
	log.Fatal(app.Listen(addr, fiber.ListenConfig{
		// Prefork spawns multiple processes that share JSON files;
		// the file store relies on in-process locking, so it stays off.
		EnablePrefork: false,
	}))
}

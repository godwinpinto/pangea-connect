package main

import (
	"log"
	"os"
	"strconv"

	pangea "github.com/godwinpinto/gofiber/middleware"
	"github.com/gofiber/fiber/v3"
)

// A Sample gin gonic rest api
func main() {
	app := fiber.New()

	//A middleware using Pangea IP Intel
	ipIntelType := os.Getenv("PANGEA_IP_INTEL_TYPE")
	if ipIntelType == "" {
		log.Fatal("Unauthorized: No Intel type present")
	}
	pangeaToken := os.Getenv("PANGEA_INTEL_TOKEN")
	if pangeaToken == "" {
		log.Fatal("Unauthorized: No Pangea Intel token present")
	}
	pangeaDomain := os.Getenv("PANGEA_DOMAIN")
	if pangeaDomain == "" {
		log.Fatal("Unauthorized: No Pangea Domain present")
	}

	scoreThresholdStr := os.Getenv("PANGEA_IP_INTEL_SCORE_THRESHOLD")
	if scoreThresholdStr == "" {
		log.Fatal("Unauthorized: No Intel type present")
	}
	score, err := strconv.Atoi(scoreThresholdStr)
	if err != nil {
		log.Fatalf("Failed to parse PANGEA_SCORE_THRESHOLD: %v", err)
	}

	app.Use(pangea.New(pangea.Config{
		IpIntelType:    ipIntelType,
		PangeaToken:    pangeaToken,
		PangeaDomain:   pangeaDomain,
		ScoreThreshold: score,
	}))

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":8080"))

}

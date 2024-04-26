// Package plugin provides middleware for handling IP intelligence checks using Pangea services.
package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
)

// ipScores is a sync.Map storing the IP addresses and their reputation scores.
var ipScores sync.Map

// PangeaIpIntel is a middleware function that checks the IP intelligence for incoming requests.
func New(config Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	// Return new handler
	return func(c fiber.Ctx) error {

		ipIntelType := cfg.IpIntelType
		if ipIntelType == "" {
			log.Fatal("Unauthorized: No Intel type present")
		}
		var ipAddress string

		// Determine the IP address based on the IP intelligence type.
		if ipIntelType == "header" {
			ip := c.Get(fiber.HeaderXForwardedFor)
			if ip != "" {
				ipAddress = ip
			}
		} else if ipIntelType == "request" {
			ip := c.Context().RemoteAddr().String()
			if ip != "" {
				ipAddress = ip
			}
		}
		if ipAddress != "" {
			isIpPresentInCache := false
			// Check if the IP address is already in the cache.
			val, _ := client.Get(context.Background(), ipAddress).Result()
			if val == "Y" {
				isIpPresentInCache = true
				return c.Status(http.StatusForbidden).SendString("Forbidden")
			} else if val == "N" {
				isIpPresentInCache = true
			}
			// If the IP address is not in the cache, fetch its reputation asynchronously.
			if !isIpPresentInCache {
				go func() {
					pangeaToken := cfg.PangeaToken
					if pangeaToken == "" {
						log.Fatal("Unauthorized: No Pangea Intel token present")
					}

					pangeaDomain := cfg.PangeaDomain
					if pangeaDomain == "" {
						log.Fatal("Unauthorized: No Pangea Domain present")
					}

					intelcli := ip_intel.New(&pangea.Config{
						Token:  pangeaToken,
						Domain: pangeaDomain,
					})

					ctx := context.Background()
					ip := ipAddress
					input := &ip_intel.IpReputationRequest{
						Ip:       ip,
						Raw:      pangea.Bool(true),
						Verbose:  pangea.Bool(true),
						Provider: "crowdstrike",
					}
					resp, err := intelcli.Reputation(ctx, input)

					if err == nil && strings.ToLower(*resp.Status) == "success" {

						// If the reputation score is above a certain threshold, flag the IP address.
						scoreThreshold := cfg.ScoreThreshold
						scoreStatus := "N"
						if resp.Result.Data.Score >= scoreThreshold {
							scoreStatus = "Y"
						} else {
							scoreStatus = "N"
						}

						err = client.Set(context.Background(), ip, scoreStatus, time.Duration(time.Minute*time.Duration(cfg.CacheTimeoutInMinutes))).Err()
						if err != nil {
							fmt.Println("Failed to set key:", err)
						}
					}
				}()
			}
		}
		return c.Next()
	}
}

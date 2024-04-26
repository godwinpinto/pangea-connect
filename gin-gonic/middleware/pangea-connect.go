// Package plugin provides middleware for handling IP intelligence checks using Pangea services.
package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
)

// ipScores is a sync.Map storing the IP addresses and their reputation scores.
var ipScores sync.Map

// PangeaIpIntel is a middleware function that checks the IP intelligence for incoming requests.
func PangeaIpIntel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipIntelType := os.Getenv("PANGEA_IP_INTEL_TYPE")
		if ipIntelType == "" {
			log.Fatal("Unauthorized: No Intel type present")
		}
		var ipAddress string

		// Determine the IP address based on the IP intelligence type.
		if ipIntelType == "header" {
			ip := c.Request.Header.Get("X-Forwarded-For")
			if ip != "" {
				ipAddress = ip
			}
		} else if ipIntelType == "request" {
			ip := c.Request.RemoteAddr
			if ip != "" {
				ipAddress = ip
			}
		}
		if ipAddress != "" {
			isIpPresentInCache := false
			// Check if the IP address is already in the cache.
			//TODO: Redis option for distributed cache
			if val, exists := ipScores.Load(ipAddress); exists {
				isIpPresentInCache = true
				if val.(bool) {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
			}

			// If the IP address is not in the cache, fetch its reputation asynchronously.
			if !isIpPresentInCache {
				go func() {
					pangeaToken := os.Getenv("PANGEA_INTEL_TOKEN")
					if pangeaToken == "" {
						log.Fatal("Unauthorized: No Pangea Intel token present")
					}

					pangeaDomain := os.Getenv("PANGEA_DOMAIN")
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
						scoreThresholdStr := os.Getenv("PANGEA_IP_INTEL_SCORE_THRESHOLD")
						if scoreThresholdStr == "" {
							log.Fatal("Unauthorized: No Intel type present")
						}

						score, err := strconv.Atoi(scoreThresholdStr)
						if err != nil {
							log.Fatalf("Failed to parse PANGEA_SCORE_THRESHOLD: %v", err)
						}
						var scoreThreshold = score
						if resp.Result.Data.Score >= scoreThreshold {
							ipScores.Store(ip, true)
						} else {
							ipScores.Store(ip, false)
						}
					}
				}()
			}
		}
		c.Next()
	}
}

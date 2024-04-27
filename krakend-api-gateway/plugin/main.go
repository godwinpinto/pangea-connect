package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
	"github.com/redis/go-redis/v9"
)

// pluginName is the plugin name
var pluginName = "krakend-pangea-connect"

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	// The config variable contains all the keys you have defined in the configuration
	// if the key doesn't exists or is not a map the plugin returns an error and the default handler
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	status, _ := config["enabled"].(bool)
	/*
	* The pangea connect plugin to fetch environment variables
	 */
	ipIntelType := os.Getenv("PANGEA_IP_INTEL_TYPE")
	if ipIntelType == "" {
		logger.Fatal("Unauthorized: No Intel type present")
	}
	pangeaToken := os.Getenv("PANGEA_INTEL_TOKEN")
	if pangeaToken == "" {
		logger.Fatal("Unauthorized: No Pangea Intel token present")
	}
	pangeaDomain := os.Getenv("PANGEA_DOMAIN")
	if pangeaDomain == "" {
		logger.Fatal("Unauthorized: No Pangea Domain present")
	}

	scoreThresholdStr := os.Getenv("PANGEA_IP_INTEL_SCORE_THRESHOLD")
	if scoreThresholdStr == "" {
		logger.Fatal("Unauthorized: No Intel type present")
	}
	score, err := strconv.Atoi(scoreThresholdStr)
	if err != nil {
		logger.Fatal("Failed to parse PANGEA_SCORE_THRESHOLD: %v", err)
	}

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		logger.Fatal("Unauthorized: No redis url present")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisTtlStr := os.Getenv("PANGEA_IP_INTEL_SCORE_THRESHOLD")
	if redisTtlStr == "" {
		logger.Fatal("Unauthorized: No Intel type present")
	}
	redisTtl, err := strconv.Atoi(redisTtlStr)
	if err != nil {
		logger.Fatal("Failed to parse PANGEA_SCORE_THRESHOLD: %v", err)
	}
	var client *redis.Client
	if status {
		client = redis.NewClient(&redis.Options{
			Addr:     redisUrl,      // Redis server address
			Password: redisPassword, // No password
			DB:       0,             // Default DB
		})
	}

	logger.Debug(fmt.Sprintf("The plugin is now hijacking the path %v", status))

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		/*
		* Business logic for Pangea Ip intel service
		*
		 */

		if !status {
			h.ServeHTTP(w, req)
			return
		}

		var ipAddress string

		// Determine the IP address based on the IP intelligence type.
		if ipIntelType == "header" {

			ip := req.Header.Get("X-Forwarded-For")
			if ip != "" {
				ipAddress = ip
			}
		} else if ipIntelType == "request" {
			ip := req.RemoteAddr
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
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "Forbidden")
				return
			} else if val == "N" {
				isIpPresentInCache = true
			}
			logger.Debug(isIpPresentInCache)
			logger.Debug(redisTtl)
			logger.Debug(score)
			// If the IP address is not in the cache, fetch its reputation asynchronously.
			if !isIpPresentInCache {
				go func() {

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
						scoreThreshold := score
						scoreStatus := "N"
						if resp.Result.Data.Score >= scoreThreshold {
							scoreStatus = "Y"
						} else {
							scoreStatus = "N"
						}

						err = client.Set(context.Background(), ip, scoreStatus, time.Duration(time.Minute*time.Duration(redisTtl))).Err()
						if err != nil {
							fmt.Println("Failed to set key:", err)
						}
					}
				}()
			}
		}
		h.ServeHTTP(w, req)
	}), nil
}

func main() {}

// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}

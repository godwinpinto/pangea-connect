package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
	"github.com/redis/go-redis/v9"
)

const PluginName = "pangea-connect"
const Version = "0.0.1"
const Priority = 1000

const FailedResponse = `{"error": "%s is required"}`

type Config struct {
	Enabled                     string `json:"enabled"`
	PangeaToken                 string `json:"pangea_token"`
	PangeaDomain                string `json:"pangea_domain"`
	RedisUrl                    string `json:"redis_url"`
	RedisPassword               string `json:"redis_password"`
	RedisTtl                    int    `json:"redis_ttl"`
	PangeaIpIntelType           string `json:"pangea_ip_intel_type"`
	PangeaIpIntelScoreThreshold int    `json:"pangea_ip_intel_score_threshold"`
	PangeaEnabled               bool   `json:"pangea_enabled"`
}

var client *redis.Client

func main() {

	errServerStart := server.StartServer(New, Version, Priority)
	if errServerStart != nil {
		log.Fatalf("Failed start %s plugin", PluginName)
	}

}

func New() interface{} {
	/*
	* The pangea connect plugin to fetch environment variables
	 */

	var conf Config
	conf.Enabled = os.Getenv("PANGEA_ENABLED")
	conf.PangeaToken = os.Getenv("PANGEA_INTEL_TOKEN")
	conf.PangeaDomain = os.Getenv("PANGEA_DOMAIN")
	conf.PangeaIpIntelType = os.Getenv("PANGEA_IP_INTEL_TYPE")
	conf.RedisUrl = os.Getenv("REDIS_URL")
	conf.RedisPassword = os.Getenv("REDIS_PASSWORD")

	scoreThresholdStr := os.Getenv("PANGEA_IP_INTEL_SCORE_THRESHOLD")
	if scoreThresholdStr == "" {
		log.Println("Unauthorized: No Intel type present")
	}
	score, err := strconv.Atoi(scoreThresholdStr)
	if err != nil {
		log.Println("Unauthorized: No score present")
	}
	conf.PangeaIpIntelScoreThreshold = score

	redisTtlStr := os.Getenv("REDIS_TTL")
	if redisTtlStr == "" {
		log.Println("Unauthorized: No Redis TTL present")
	}
	redisTtl, err := strconv.Atoi(redisTtlStr)
	if err != nil {
		log.Println("Unauthorized: No Redis TTL present")
	}
	conf.RedisTtl = redisTtl
	statusStr := os.Getenv("PANGEA_ENABLED")
	if statusStr == "" {
		log.Fatal("Unauthorized: No Intel type present")
	}

	conf.PangeaEnabled, err = strconv.ParseBool(statusStr)
	if err != nil {
		log.Fatalf("Error parsing PANGEA_ENABLED value: %v", err)
	}

	if conf.PangeaEnabled {
		log.Println("Start redis client")
		client = redis.NewClient(&redis.Options{
			Addr:     conf.RedisUrl,
			Password: conf.RedisPassword,
			DB:       0,
		})
	}

	log.Printf("Config: %+v", conf)
	return &Config{
		Enabled:                     conf.Enabled,
		PangeaToken:                 conf.PangeaToken,
		PangeaDomain:                conf.PangeaDomain,
		RedisUrl:                    conf.RedisUrl,
		RedisPassword:               conf.RedisPassword,
		RedisTtl:                    conf.RedisTtl,
		PangeaIpIntelType:           conf.PangeaIpIntelType,
		PangeaIpIntelScoreThreshold: conf.PangeaIpIntelScoreThreshold,
		PangeaEnabled:               conf.PangeaEnabled,
	}
}

func (conf *Config) Access(kong *pdk.PDK) {
	log.Println(conf.PangeaToken)

	var ipAddress string

	if conf.PangeaIpIntelType == "header" {
		ip, _ := kong.Request.GetHeader("X-Forwarded-For")
		//ip := req.Header.Get("X-Forwarded-For")
		if ip != "" {
			ipAddress = ip
		}
	} else if conf.PangeaIpIntelType == "request" {
		ip, _ := kong.Client.GetIp()
		//ip := req.RemoteAddr
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
			kong.Response.Exit(400, "Forbidden", nil)
			return
		} else if val == "N" {
			isIpPresentInCache = true
		}
		// If the IP address is not in the cache, fetch its reputation asynchronously.
		if !isIpPresentInCache {
			go func() {

				intelcli := ip_intel.New(&pangea.Config{
					Token:  conf.PangeaToken,
					Domain: conf.PangeaDomain,
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
					scoreThreshold := conf.PangeaIpIntelScoreThreshold
					scoreStatus := "N"
					if resp.Result.Data.Score >= scoreThreshold {
						scoreStatus = "Y"
					} else {
						scoreStatus = "N"
					}

					err = client.Set(context.Background(), ip, scoreStatus, time.Duration(time.Minute*time.Duration(conf.RedisTtl))).Err()
					if err != nil {
						fmt.Println("Failed to set key:", err)
					}
				}
			}()
		}
	}
}

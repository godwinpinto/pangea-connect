package middleware

// Config defines the config for middleware.
type Config struct {
	IpIntelType           string
	PangeaToken           string
	PangeaDomain          string
	ScoreThreshold        int
	RedisUrl              string
	RedisPassword         string
	CacheTimeoutInMinutes int
	// // Next defines a function to skip this middleware when returned true.
	// //
	// // Optional. Default: nil
	// Next func(c fiber.Ctx) bool

	// // Users defines the allowed credentials
	// //
	// // Required. Default: map[string]string{}
	// Users map[string]string

	// // Realm is a string to define realm attribute of BasicAuth.
	// // the realm identifies the system to authenticate against
	// // and can be used by clients to save credentials
	// //
	// // Optional. Default: "Restricted".
	// Realm string

	// // Authorizer defines a function you can pass
	// // to check the credentials however you want.
	// // It will be called with a username and password
	// // and is expected to return true or false to indicate
	// // that the credentials were approved or not.
	// //
	// // Optional. Default: nil.
	// Authorizer func(string, string) bool

	// // Unauthorized defines the response body for unauthorized responses.
	// // By default it will return with a 401 Unauthorized and the correct WWW-Auth header
	// //
	// // Optional. Default: nil
	// Unauthorized fiber.Handler
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	IpIntelType:           "none",
	PangeaToken:           "",
	PangeaDomain:          "",
	ScoreThreshold:        100,
	RedisUrl:              "localhost:6379",
	RedisPassword:         "",
	CacheTimeoutInMinutes: 15,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.IpIntelType == "" {
		cfg.IpIntelType = "none"
	}
	if cfg.PangeaToken == "" {
		cfg.PangeaToken = ConfigDefault.PangeaToken
	}
	if cfg.PangeaDomain == "" {
		cfg.PangeaDomain = ConfigDefault.PangeaDomain
	}
	if cfg.ScoreThreshold == 0 {
		cfg.ScoreThreshold = ConfigDefault.ScoreThreshold
	}
	if cfg.RedisUrl == "" {
		cfg.RedisUrl = ConfigDefault.RedisUrl
	}
	if cfg.RedisPassword == "" {
		cfg.RedisPassword = ConfigDefault.RedisPassword
	}
	if cfg.CacheTimeoutInMinutes == -1 {
		cfg.CacheTimeoutInMinutes = ConfigDefault.CacheTimeoutInMinutes
	}
	return cfg
}

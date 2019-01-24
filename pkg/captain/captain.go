package captain

const (
	EnvEndpoint = "CAPTAIN_ENDPOINT"
	EnvPort     = "CAPTAIN_PORT"
	DefaultPort = 30051

	// helm tiller
	EnvTillerEndpoint              = "TILLER_ENDPOINT"
	EnvTillerUsername              = "TILLER_USERNAME"
	EnvTillerPassword              = "TILLER_PASSWORD"
	EnvTillerAccount               = "TILLER_ACCOUNT"
	EnvTillerSkipSslValidation     = "TILLER_SKIP_SSL_VALIDATION"
	DefaultTillerUsername          = "admin"
	DefaultTillerPassword          = "admin"
	DefaultTillerAccount           = "mycluster Account"
	DefaultTillerSkipSslValidation = true

	// hub.softleader.com.tw
	EnvRegistryAuthUsername     = "REGISTRY_AUTH_USERNAME"
	EnvRegistryAuthPassword     = "REGISTRY_AUTH_PASSWORD"
	DefaultRegistryAuthUsername = "client"
	DefaultRegistryAuthPassword = "poweredbysoftleader"
)

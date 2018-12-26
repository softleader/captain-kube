package captain

const (
	EnvPort          = "CAPTAIN_PORT"
	EnvK8sVendor     = "CAPTAIN_K8S_VENDOR"
	DefaultPort      = 30051
	DefaultK8sVendor = "icp"

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
	DefaultRegistryAuthUsername = "dev"
	DefaultRegistryAuthPassword = "sleader"
)

package captain

const (
	// EnvEndpoint key to specify captain endpoint
	EnvEndpoint = "CAPTAIN_ENDPOINT"
	// EnvPort key to specify captain port
	EnvPort = "CAPTAIN_PORT"
	// DefaultPort captain default port
	DefaultPort = 30051

	// EnvTillerEndpoint key to specify helm tiller endpoint
	EnvTillerEndpoint = "TILLER_ENDPOINT"
	// EnvTillerUsername key to specify helm tiller username
	EnvTillerUsername = "TILLER_USERNAME"
	// EnvTillerPassword key to specify helm tiller password
	EnvTillerPassword = "TILLER_PASSWORD"
	// EnvTillerAccount key to specify helm tiller account
	EnvTillerAccount = "TILLER_ACCOUNT"
	// EnvTillerSkipSslValidation key to specify helm tiller should skip ssl validation
	EnvTillerSkipSslValidation = "TILLER_SKIP_SSL_VALIDATION"
	// DefaultTillerUsername default helm tiller username
	DefaultTillerUsername = "admin"
	// DefaultTillerPassword default helm tiller password
	DefaultTillerPassword = "admin"
	// DefaultTillerAccount default helm tiller account
	DefaultTillerAccount = "mycluster Account"
	// DefaultTillerSkipSslValidation default helm tiller should skip ssl validation
	DefaultTillerSkipSslValidation = true

	// EnvRegistryAuthUsername key to specify docker registry username
	EnvRegistryAuthUsername = "REGISTRY_AUTH_USERNAME"
	// EnvRegistryAuthPassword key to specify docker registry password
	EnvRegistryAuthPassword = "REGISTRY_AUTH_PASSWORD"
	// DefaultRegistryAuthUsername default docker registry username
	DefaultRegistryAuthUsername = "client"
	// DefaultRegistryAuthPassword default docker registry password
	DefaultRegistryAuthPassword = "poweredbysoftleader"
)

package captain

const (
	EnvTillerEndpoint          = "TILLER_ENDPOINT"
	EnvTillerUsername          = "TILLER_USERNAME"
	EnvTillerPassword          = "TILLER_PASSWORD"
	EnvTillerAccount           = "TILLER_ACCOUNT"
	EnvTillerSkipSslValidation = "TILLER_SKIP_SSL_VALIDATION"

	DefaultTillerEndpoint          = "k8s.master.vip:8443"
	DefaultTillerUsername          = "admin"
	DefaultTillerPassword          = "admin"
	DefaultTillerAccount           = "icp Account"
	DefaultTillerSkipSslValidation = true
)

package handler

type Handlers struct {
	AuthHandler     *AuthHandler
	ApiKeyHandler   *APIKeyHandler
	GatewayHandler  *GatewayHandler
	LogHandler      *LogHandler
	ProviderHandler *ProviderHandler
	TenantHandler   *TenantHandler
	MetricsHandler  *MetricsHandler
	StatsHandler    *StatsHandler
}

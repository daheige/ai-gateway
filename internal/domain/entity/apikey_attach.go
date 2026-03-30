package entity

// APIKeyCreateEntity apikey创建实体
type APIKeyCreateEntity struct {
	TenantID          uint
	ProviderID        uint
	Name              string
	Prefix            string
	RateLimitPerSec   int
	RateLimit         int
	MonthlyTokenLimit int
	TotalTokenQuota   int
}

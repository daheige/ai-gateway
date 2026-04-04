package entity

// DailyStat 每日统计
type DailyStat struct {
	Date   string `json:"date"`
	Tokens int64  `json:"tokens"`
}

// KeyStat 按Key统计
type KeyStat struct {
	APIKeyID  uint   `json:"api_key_id"`
	KeyPrefix string `json:"key_prefix"`
	Tokens    int64  `json:"tokens"`
}

// Overview 总览
type StatsOverview struct {
	TotalTokensToday int64 `json:"total_tokens_today"`
	TotalTokensMonth int64 `json:"total_tokens_month"`
	TotalRequests    int64 `json:"total_requests"`
	TotalRequestsDay int64 `json:"total_requests_today"`
	ActiveKeys       int64 `json:"active_keys"`
	ActiveProviders  int64 `json:"active_providers"`
}

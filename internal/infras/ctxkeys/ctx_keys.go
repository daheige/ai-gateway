package ctxkeys

// CtxKey ctx key struct.
type CtxKey struct {
	Name string
}

// String CtxKey string.
func (c CtxKey) String() string {
	return c.Name
}

var (
	// XRequestID request_id
	XRequestID = CtxKey{"x-request-id"}

	// ClientIP request client ip
	ClientIP = CtxKey{"client_ip"}

	// RequestMethod request method
	RequestMethod = CtxKey{"request_method"}

	// RequestURI request uri
	RequestURI = CtxKey{"request_uri"}

	// UserAgent request ua
	UserAgent = CtxKey{"request_ua"}

	// CurHostname current hostname
	CurHostname = CtxKey{"hostname"}

	// APIKeyInfo apikey info
	APIKeyInfo = CtxKey{"api_key_info"}

	APIKey = CtxKey{"api_key"}

	TenantID = CtxKey{"tenant_id"}
)

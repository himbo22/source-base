package constraints

type contextKey string

const (
	RequestIDKey  contextKey = "request_id"
	ClaimsKey     contextKey = "claims"
	PublicIDKey   contextKey = "public_id"
	UserStatusKey contextKey = "user_status"
	TxKey         contextKey = "tx"
)

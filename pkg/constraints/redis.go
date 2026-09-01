package constraints

import "fmt"

const (
	// Domains
	DomainAuth = "auth"
	DomainUser = "user"
	DomainSys  = "sys"
)

// RedisKeyBuilder is a utility namespace for generating standard Redis keys
type RedisKeyBuilder struct{}

var RedisKey = RedisKeyBuilder{}

// --- Auth Domain Keys ---

// ATBlacklist generates the key for blacklisted Access Tokens (JTI)
func (RedisKeyBuilder) ATBlacklist(jti string) string {
	return fmt.Sprintf("%s:blacklist:at:%s", DomainAuth, jti)
}

// RefreshToken generates the key for Refresh Tokens
func (RedisKeyBuilder) RefreshToken(rtID string) string {
	return fmt.Sprintf("%s:rt:%s", DomainAuth, rtID)
}

// UserStatus generates the key for user real-time status/restrictions
func (RedisKeyBuilder) UserStatus(identifier string) string {
	return fmt.Sprintf("%s:user:status:%s", DomainAuth, identifier)
}

// --- Future Domains ---
// func (RedisKeyBuilder) UserProfileCache(userID string) string {
// 	return fmt.Sprintf("%s:profile:%s", DomainUser, userID)
// }

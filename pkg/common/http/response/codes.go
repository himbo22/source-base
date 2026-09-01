package response

import "net/http"

const (
	// CodeSuccess Success codes (20000-29999)
	CodeSuccess   = 20000 // Success
	CodeCreated   = 20001 // Resource created successfully
	CodeUpdated   = 20002 // Resource updated successfully
	CodeDeleted   = 20003 // Resource deleted successfully
	CodeRetrieved = 20004 // Resource retrieved successfully

	// CodeParamInvalid Client error codes (40000-49999)
	CodeParamInvalid     = 40000 // Invalid parameters
	CodeValidationFailed = 40001 // Validation failed
	CodeBadRequest       = 40002 // Bad request
	CodeInvalidID        = 40003 // Invalid ID format

	// CodeUnauthorized Authentication/Authorization errors (41000-41999)
	CodeUnauthorized    = 41000 // Unauthorized
	CodeInvalidToken    = 41001 // Invalid token
	CodeTokenExpired    = 41002 // Token expired
	CodeInvalidPassword = 41003 // Invalid password
	CodeAccountNotFound = 41004 // Account not found
	CodeTokenRevoked    = 41005 // Token revoked

	// CodeForbidden Forbidden & Restriction errors (43000-43999)
	CodeForbidden        = 43000 // Forbidden
	CodeUserSuspended    = 43001 // User account suspended
	CodeUserRestricted   = 43002 // User account restricted
	CodeActionRestricted = 43003 // User action restricted

	// CodeNotFound Not found errors (44000-44999)
	CodeNotFound = 44000 // Resource not found

	// CodeMethodNotAllowed Method not allowed
	CodeMethodNotAllowed = 45000

	// CodeConflict Conflict errors (49000-49999)
	CodeConflict = 49000 // Conflict

	// CodeInternalServer Server error codes (50000-59999)
	CodeInternalServer = 50000 // Internal server error
	CodeDatabaseError  = 50001 // Database error
	CodeMongoDBError   = 50002 // MongoDB error
	CodeRedisError     = 50003 // Redis error
)

func MapHTTPCodeToAppCode(httpCode int) int {
	switch httpCode {
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusMethodNotAllowed:
		return CodeMethodNotAllowed
	default:
		return CodeInvalidID
	}
}

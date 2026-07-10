package response

import "net/http"

type CodeInfo struct {
	HTTPStatus int
	Message    string
}

// Msg maps codes to default messages
var Msg = map[int]CodeInfo{
	// Success
	CodeSuccess:   {http.StatusOK, "Success"},
	CodeCreated:   {http.StatusCreated, "Resource created successfully"},
	CodeUpdated:   {http.StatusOK, "Resource updated successfully"},
	CodeDeleted:   {http.StatusOK, "Resource deleted successfully"},
	CodeRetrieved: {http.StatusOK, "Resource retrieved successfully"},

	CodeParamInvalid:     {http.StatusBadRequest, "Invalid parameters"},
	CodeValidationFailed: {http.StatusBadRequest, "Validation failed"},
	CodeBadRequest:       {http.StatusBadRequest, "Bad request"},
	CodeInvalidID:        {http.StatusBadRequest, "Invalid ID format"},

	CodeUnauthorized:    {http.StatusUnauthorized, "Unauthorized"},
	CodeInvalidToken:    {http.StatusUnauthorized, "Invalid token"},
	CodeTokenExpired:    {http.StatusUnauthorized, "Token expired"},
	CodeInvalidPassword: {http.StatusUnauthorized, "Invalid password"},
	CodeAccountNotFound: {http.StatusUnauthorized, "Account not found"},
	CodeForbidden:       {http.StatusForbidden, "Forbidden"},

	CodeNotFound: {http.StatusNotFound, "Resource not found"},

	CodeConflict: {http.StatusConflict, "Conflict"},

	CodeInternalServer: {http.StatusInternalServerError, "Internal server error"},
	CodeDatabaseError:  {http.StatusInternalServerError, "Database error"},
	CodeMongoDBError:   {http.StatusInternalServerError, "MongoDB error"},
	CodeRedisError:     {http.StatusInternalServerError, "Redis error"},
}

func Get(code int) CodeInfo {
	if info, ok := Msg[code]; ok {
		return info
	}

	return CodeInfo{
		HTTPStatus: http.StatusInternalServerError,
		Message:    "Internal server error",
	}
}

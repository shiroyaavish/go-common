package common

import (
	"github.com/IntelXLabs-LLC/go-common/config"
	"github.com/IntelXLabs-LLC/go-common/errors/api_errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
)

// AccessLevel represents the access level of a user or service in the system.
type AccessLevel byte

const (

	// AccessLevelPublic is a constant of type AccessLevel with a value of 0.
	AccessLevelPublic = AccessLevel(0)
	// AccessLevelUser is a constant of type AccessLevel with a value of 1.
	AccessLevelUser = AccessLevel(1)
	// AccessLevelService is a constant of type AccessLevel with a value of 2.
	AccessLevelService = AccessLevel(2)
	// AccessLevelAdmin is a constant of type AccessLevel with a value of 3.
	AccessLevelAdmin = AccessLevel(3)
	// AccessLevelSearch is a constant of type AccessLevel with a value of 4.
	AccessLevelSearch = AccessLevel(4)
)

// ToString returns a string representation of the AccessLevel value.
// It maps the AccessLevel enum values to their corresponding string representation.
// If the AccessLevel value is not one of the defined enum values, the string "unknown" is returned.
func (a AccessLevel) ToString() string {
	switch a {
	case AccessLevelPublic:
		return "public"
	case AccessLevelUser:
		return "user"
	case AccessLevelService:
		return "service"
	case AccessLevelAdmin:
		return "admin"
	case AccessLevelSearch:
		return "search"
	}
	return "unknown"
}

// CheckAccess checks the access level and returns the necessary information based on the access level provided.
// If the access level is AccessLevelPublic, it returns nil, nil, nil.
// If the access level is AccessLevel
func (a AccessLevel) CheckAccess(c *fiber.Ctx) (*uuid.UUID, []string, *api_errors.Error) {
	// Checks access according to the values
	switch a {
	case AccessLevelPublic:
		return nil, nil, nil
	case AccessLevelService:
		return checkServiceAccess(c.Get("Api-Key", ""))
	case AccessLevelSearch:
		_, _, err := checkServiceAccess(c.Get("X-Api-Key", ""))
		if err == nil {
			return nil, nil, nil
		}
		fallthrough
	case AccessLevelUser, AccessLevelAdmin:
		authHeader := c.Get("Authorization", "")
		if authHeader == "" {
			return nil, nil, api_errors.ErrUnauthorized
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return nil, nil, api_errors.ErrUnauthorized
		}
		return checkUserAndAdminAccess(authHeaderParts[1], a)
	default:
		return nil, nil, api_errors.ErrUnauthorized
	}
}

// checkUserAndAdminAccess checks the user and admin access based on the provided token and access level.
// It parses the JWT token and verifies the access level.
// If the access level is AccessLevelAdmin, it uses the admin secret from the config.
// If the access level is not AccessLevelAdmin, it uses the user secret from the config.
// It validates the token claims and checks for the access level.
// If the access level is AccessLevelSearch, it checks if the "search" access is allowed.
// If the access level is not allowed, it returns an unauthorized error.
// It retrieves the user ID from the token claims and parses it into a UUID.
// It retrieves the permissions from the token claims and converts them into a string slice.
// It returns the access ID, access permissions, and nil error if successful.
// If any error occurs during parsing, validation, or retrieving claims, it returns nil access ID, nil access permissions, and an unauthorized error.
func checkUserAndAdminAccess(tokenStr string, accessLevel AccessLevel) (*uuid.UUID, []string, *api_errors.Error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, api_errors.ErrUnauthorized
		}

		if accessLevel == AccessLevelAdmin {
			return []byte(config.GetAdminSecret()), nil
		}
		return []byte(config.GetUserSecret()), nil
	})
	if err != nil {
		return nil, nil, api_errors.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, api_errors.ErrUnauthorized
	}

	if accessLevel == AccessLevelSearch {
		access, ok := claims["access"].(map[string]bool)
		if !ok {
			return nil, nil, api_errors.ErrUnauthorized
		}
		if allowed, ok := access["search"]; ok && allowed {
			return nil, nil, nil
		}
		return nil, nil, api_errors.ErrUnauthorized
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, nil, api_errors.ErrUnauthorized
	}
	accessId, err := uuid.Parse(id)
	if err != nil {
		return nil, nil, api_errors.ErrUnauthorized
	}
	accessPermissions := make([]string, 0)
	permissions, ok := claims["permissions"].([]any)
	if ok {
		for _, v := range permissions {
			accessPermissions = append(accessPermissions, v.(string))
		}
	}

	return &accessId, accessPermissions, nil
}

// checkServiceAccess checks the service access based on the provided API key.
// If the API key matches the service secret from the config, it returns nil access ID, nil access permissions, and nil error.
// If the API key does not match the service secret, it returns nil access ID, nil access permissions, and an unauthorized error.
func checkServiceAccess(apiKey string) (*uuid.UUID, []string, *api_errors.Error) {
	if apiKey == config.GetServiceSecret() {
		return nil, nil, nil
	}
	return nil, nil, api_errors.ErrUnauthorized
}

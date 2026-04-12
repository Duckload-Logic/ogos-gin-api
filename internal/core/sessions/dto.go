package sessions

import (
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// JTIDTO encapsulates the JWT ID (jti) and provides helper methods
// for generating Redis keys to ensure consistency and maintainability.
type JTIDTO struct {
	Value string
}

// NewJTI creates a new JTIDTO instance.
func NewJTI(value string) JTIDTO {
	return JTIDTO{Value: value}
}

// ToSessionKey returns the Redis key for the primary session data.
func (j JTIDTO) ToSessionKey() string {
	return fmt.Sprintf("%s%s", constants.RedisSessionKeyPrefix, j.Value)
}

// ToIDPRefreshKey returns the Redis key for the linked IDP refresh
// token.
func (j JTIDTO) ToIDPRefreshKey() string {
	return fmt.Sprintf("%s%s", constants.RedisIDPRefreshKeyPrefix, j.Value)
}

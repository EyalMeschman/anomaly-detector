package validator

import (
	"regexp"

	"anomaly_detector/models"
)

var (
	// Date format: dd-mm-yyyy
	dateRegex = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])-(0[1-9]|1[0-2])-\d{4}$`)
	// Email format: simplified RFC 5321
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// Auth-Token format: Bearer <token>
	authTokenRegex = regexp.MustCompile(`^Bearer [a-zA-Z0-9]+$`)
)

func validateType(value any, typeName models.ParamType) bool {
	switch v := value.(type) {
	case string:
		return validateStringType(v, typeName)

	case float64:
		// JSON unmarshaling converts numbers to float64
		if typeName != models.TypeInt {
			return false
		}
		// Ensure it's actually an integer
		return v == float64(int(v))

	case int, int32, int64:
		return typeName == models.TypeInt

	case bool:
		return typeName == models.TypeBoolean

	case []any, []map[string]any:
		return typeName == models.TypeList

	default:
		return false
	}
}

func validateStringType(value string, typeName models.ParamType) bool {
	switch typeName {
	case models.TypeString:
		return true

	case models.TypeDate:
		return dateRegex.MatchString(value)

	case models.TypeEmail:
		return emailRegex.MatchString(value)

	case models.TypeUUID:
		// We explicitly require UUIDs with dashes (-). If we didn’t, I’d use the official github.com/google/uuid package
		return uuidRegex.MatchString(value)

	case models.TypeAuthToken:
		return authTokenRegex.MatchString(value)

	default:
		return false
	}
}

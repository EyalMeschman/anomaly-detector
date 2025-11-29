package validator

import (
	"fmt"
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

// validateType validates a value against a type name
func validateType(value any, typeName string) error {
	switch typeName {
	case models.TypeString:
		_, err := validateString(value)
		return err
	case models.TypeInt:
		return validateInt(value)
	case models.TypeBoolean:
		return validateBoolean(value)
	case models.TypeList:
		return validateList(value)
	case models.TypeDate:
		return validateDate(value)
	case models.TypeEmail:
		return validateEmail(value)
	case models.TypeUUID:
		return validateUUID(value)
	case models.TypeAuthToken:
		return validateAuthToken(value)
	default:
		return fmt.Errorf("unknown type: %s", typeName)
	}
}

func validateInt(value any) error {
	switch v := value.(type) {
	case float64:
		// JSON unmarshaling converts numbers to float64
		if v != float64(int(v)) {
			return fmt.Errorf("value %v is a float, not an integer", v)
		}

		return nil
	case int, int32, int64:
		// Only happens if developer manually sets ints
		return nil
	default:
		return fmt.Errorf("expected Int, got %T", value)
	}
}

func validateString(value any) (string, error) {
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected String, got %T", value)
	}

	return strValue, nil
}

func validateBoolean(value any) error {
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("expected Boolean, got %T", value)
	}

	return nil
}

func validateList(value any) error {
	switch value.(type) {
	case []any, []map[string]any:
		return nil
	default:
		return fmt.Errorf("expected List, got %T", value)
	}
}

func validateDate(value any) error {
	strValue, err := validateString(value)
	if err != nil {
		return fmt.Errorf("expected Date string. err: %v", err)
	}

	if !dateRegex.MatchString(strValue) {
		return fmt.Errorf("invalid date format: %s", strValue)
	}

	return nil
}

func validateEmail(value any) error {
	strValue, err := validateString(value)
	if err != nil {
		return fmt.Errorf("expected Email string. err: %v", err)
	}

	if !emailRegex.MatchString(strValue) {
		return fmt.Errorf("invalid email format: %s", strValue)
	}

	return nil
}

func validateUUID(value any) error {
	strValue, err := validateString(value)
	if err != nil {
		return fmt.Errorf("expected UUID string. err: %v", err)
	}

	if !uuidRegex.MatchString(strValue) {
		return fmt.Errorf("invalid UUID format: %s", strValue)
	}

	return nil
}

func validateAuthToken(value any) error {
	strValue, err := validateString(value)
	if err != nil {
		return fmt.Errorf("expected Auth-Token string. err: %v", err)
	}

	if !authTokenRegex.MatchString(strValue) {
		return fmt.Errorf("invalid Auth-Token format: %s", strValue)
	}

	return nil
}

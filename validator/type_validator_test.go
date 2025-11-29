package validator

import (
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

type typeTestCase struct {
	name       string
	inputValue any
	isValid    bool
}

func runTypeTests(t *testing.T, typeName models.ParamType, testCases []typeTestCase) {
	t.Run(string(typeName), func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := validateType(tc.inputValue, typeName)
				assert.Equal(t, tc.isValid, result)
			})
		}
	})
}

func TestValidateType(t *testing.T) {
	runTypeTests(t, models.TypeInt, []typeTestCase{
		{name: "valid int from json", inputValue: float64(123), isValid: true},
		{name: "valid int whole number", inputValue: float64(100.0), isValid: true},
		{name: "invalid float", inputValue: float64(123.5), isValid: false},
		{name: "valid go int", inputValue: 123, isValid: true},
		{name: "invalid type", inputValue: "123", isValid: false},
	})

	runTypeTests(t, models.TypeString, []typeTestCase{
		{name: "valid string", inputValue: "test", isValid: true},
		{name: "empty string valid", inputValue: "", isValid: true},
		{name: "invalid type", inputValue: 123, isValid: false},
	})

	runTypeTests(t, models.TypeBoolean, []typeTestCase{
		{name: "valid true", inputValue: true, isValid: true},
		{name: "valid false", inputValue: false, isValid: true},
		{name: "invalid type", inputValue: "true", isValid: false},
	})

	runTypeTests(t, models.TypeList, []typeTestCase{
		{name: "valid empty list", inputValue: []any{}, isValid: true},
		{name: "valid list with items", inputValue: []any{1, 2, 3}, isValid: true},
		{name: "invalid type", inputValue: "not a list", isValid: false},
	})

	runTypeTests(t, models.TypeDate, []typeTestCase{
		{name: "valid date", inputValue: "12-01-2022", isValid: true},
		{name: "valid date", inputValue: "31-12-2023", isValid: true},
		{name: "invalid day", inputValue: "67-12-2023", isValid: false},
		{name: "invalid month", inputValue: "14-45-2023", isValid: false},
		{name: "invalid format", inputValue: "2022-01-12", isValid: false},
		{name: "invalid type", inputValue: 123, isValid: false},
	})

	runTypeTests(t, models.TypeEmail, []typeTestCase{
		{name: "valid email", inputValue: "test@example.com", isValid: true},
		{name: "valid with subdomain", inputValue: "user@mail.example.co.uk", isValid: true},
		{name: "invalid no @", inputValue: "notanemail", isValid: false},
		{name: "invalid type", inputValue: 123, isValid: false},
	})

	runTypeTests(t, models.TypeUUID, []typeTestCase{
		{name: "valid uuid", inputValue: "46da6390-7c78-4a1c-9efa-7c0396067ce4", isValid: true},
		{name: "valid uuid lowercase", inputValue: "550e8400-e29b-41d4-a716-446655440000", isValid: true},
		{name: "invalid too short", inputValue: "46da6390-7c78-4a1c-9efa", isValid: false},
		{name: "invalid no dashes", inputValue: "46da63907c784a1c9efa7c0396067ce4", isValid: false},
		{name: "invalid type", inputValue: 123, isValid: false},
	})

	runTypeTests(t, models.TypeAuthToken, []typeTestCase{
		{name: "valid token", inputValue: "Bearer abc123", isValid: true},
		{name: "valid long token", inputValue: "Bearer ebb3cbbe938c4776bd22a4ec2ea8b2ca", isValid: true},
		{name: "invalid no Bearer", inputValue: "abc123", isValid: false},
		{name: "invalid lowercase bearer", inputValue: "bearer abc123", isValid: false},
		{name: "invalid special chars", inputValue: "Bearer abc-123", isValid: false},
		{name: "invalid type", inputValue: 123, isValid: false},
	})

	t.Run("UnknownType", func(t *testing.T) {
		result := validateType("value", "UnknownType")
		assert.False(t, result)
	})
}

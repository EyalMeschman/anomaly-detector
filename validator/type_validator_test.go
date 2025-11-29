package validator

import (
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

type typeTestCase struct {
	name       string
	inputValue any
	isErr      bool
}

func runTypeTests(t *testing.T, typeName string, testCases []typeTestCase) {
	t.Run(typeName, func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := validateType(tc.inputValue, typeName)
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestValidateType(t *testing.T) {
	runTypeTests(t, models.TypeInt, []typeTestCase{
		{name: "valid int from json", inputValue: float64(123), isErr: false},
		{name: "valid int whole number", inputValue: float64(100.0), isErr: false},
		{name: "invalid float", inputValue: float64(123.5), isErr: true},
		{name: "valid go int", inputValue: 123, isErr: false},
		{name: "invalid type", inputValue: "123", isErr: true},
	})

	runTypeTests(t, models.TypeString, []typeTestCase{
		{name: "valid string", inputValue: "test", isErr: false},
		{name: "empty string valid", inputValue: "", isErr: false},
		{name: "invalid type", inputValue: 123, isErr: true},
	})

	runTypeTests(t, models.TypeBoolean, []typeTestCase{
		{name: "valid true", inputValue: true, isErr: false},
		{name: "valid false", inputValue: false, isErr: false},
		{name: "invalid type", inputValue: "true", isErr: true},
	})

	runTypeTests(t, models.TypeList, []typeTestCase{
		{name: "valid empty list", inputValue: []any{}, isErr: false},
		{name: "valid list with items", inputValue: []any{1, 2, 3}, isErr: false},
		{name: "invalid type", inputValue: "not a list", isErr: true},
	})

	runTypeTests(t, models.TypeDate, []typeTestCase{
		{name: "valid date", inputValue: "12-01-2022", isErr: false},
		{name: "valid date", inputValue: "31-12-2023", isErr: false},
		{name: "invalid format", inputValue: "2022-01-12", isErr: true},
		{name: "invalid type", inputValue: 123, isErr: true},
	})

	runTypeTests(t, models.TypeEmail, []typeTestCase{
		{name: "valid email", inputValue: "test@example.com", isErr: false},
		{name: "valid with subdomain", inputValue: "user@mail.example.co.uk", isErr: false},
		{name: "invalid no @", inputValue: "notanemail", isErr: true},
		{name: "invalid type", inputValue: 123, isErr: true},
	})

	runTypeTests(t, models.TypeUUID, []typeTestCase{
		{name: "valid uuid", inputValue: "46da6390-7c78-4a1c-9efa-7c0396067ce4", isErr: false},
		{name: "valid uuid lowercase", inputValue: "550e8400-e29b-41d4-a716-446655440000", isErr: false},
		{name: "invalid too short", inputValue: "46da6390-7c78-4a1c-9efa", isErr: true},
		{name: "invalid no dashes", inputValue: "46da63907c784a1c9efa7c0396067ce4", isErr: true},
		{name: "invalid type", inputValue: 123, isErr: true},
	})

	runTypeTests(t, models.TypeAuthToken, []typeTestCase{
		{name: "valid token", inputValue: "Bearer abc123", isErr: false},
		{name: "valid long token", inputValue: "Bearer ebb3cbbe938c4776bd22a4ec2ea8b2ca", isErr: false},
		{name: "invalid no Bearer", inputValue: "abc123", isErr: true},
		{name: "invalid lowercase bearer", inputValue: "bearer abc123", isErr: true},
		{name: "invalid special chars", inputValue: "Bearer abc-123", isErr: true},
		{name: "invalid type", inputValue: 123, isErr: true},
	})

	t.Run("UnknownType", func(t *testing.T) {
		err := validateType("value", "UnknownType")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown type")
	})
}

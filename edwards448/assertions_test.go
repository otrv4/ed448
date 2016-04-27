package edwards448

import (
	"testing"
	"reflect"
	"math/big"
)

func assert_true(t *testing.T, actual bool) {
	assert_equals(t, actual, true)
}

func assert_equals(t *testing.T, actual, expected interface{}) {
	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(expected)
	if !check_type_equality(actualType, expectedType) {
		t.Errorf("The expected type for the value is %v but intead got %v.",
			expectedType, actualType)
	}

	if !check_value_equality(actual, expected) {
		t.Errorf("The expected value was <%v> but instead got <%v>.",
			actual, expected)
	}
}

func check_type_equality(actualType, expectedType interface{}) bool {
	return expectedType == actualType
}

func check_value_equality(actual, expected interface{}) bool {
	switch actual.(type) {
	case *big.Rat:
		return expected.(*big.Rat).Cmp(actual.(*big.Rat)) == 0
	default:
		return actual == expected
	}
}


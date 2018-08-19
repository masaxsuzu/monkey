package object

import (
	"testing"
)

func TestObjects(t *testing.T) {
	tests := []struct {
		input              Object
		expectedObjectType ObjectType
	}{
		{input: &Boolean{}, expectedObjectType: BOOLEAN_OBJ},
		{input: &Integer{}, expectedObjectType: INTEGER_OBJ},
		{input: &String{}, expectedObjectType: STRING_OBJ},
		{input: &ReturnValue{}, expectedObjectType: RETURN_VALUE_OBJ},
		{input: &Error{}, expectedObjectType: ERROR_OBJ},
		{input: &Null{}, expectedObjectType: NULL_OBJ},
	}

	for _, tt := range tests {
		testObject(t, tt.input, tt.expectedObjectType)
	}
}

func testObject(t *testing.T, obj Object, expected ObjectType) {
	if obj.Type() != expected {
		t.Fatalf("obj.Type() is different from %T. got=%T", expected, obj.Type())
	}
}

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
		{input: &Quote{}, expectedObjectType: QUOTE_OBJ},
	}

	for _, tt := range tests {
		testObject(t, tt.input, tt.expectedObjectType)
	}
}

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("trues do not have same hash key")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("falses do not have same hash key")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("true has same hash key as false")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two1 := &Integer{Value: 2}
	two2 := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("integers with twoerent content have same hash keys")
	}
}

func TestHashType(t *testing.T) {
	hash := Hash{}
	if hash.Type() != "HASH" {
		t.Errorf("Hash is not hash")
	}
}

func testObject(t *testing.T, obj Object, expected ObjectType) {
	if obj.Type() != expected {
		t.Fatalf("obj.Type() is different from %T. got=%T", expected, obj.Type())
	}
}

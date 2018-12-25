package vm

import (
	"fmt"
	"github.com/masa-suzu/monkey/ast"
	"github.com/masa-suzu/monkey/compiler"
	"github.com/masa-suzu/monkey/lexer"
	"github.com/masa-suzu/monkey/object"
	"github.com/masa-suzu/monkey/parser"
	"testing"
)

type testCase struct {
	in   string
	want interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		{"1", 1},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"2 * 2", 4},
		{"1 / 2", 0},
		{"-1", -1},
		{"-1 * 5", -5},
	}
	testRun(t, tests)
}
func TestIntegerArithmeticError(t *testing.T) {
	tests := []testCase{
		{"1 / 0", fmt.Errorf("integer divide by zero")},
	}
	testRunWithError(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []testCase{
		{"true", true},
		{"false", false},
		{"1 == 1", true},
		{"true == false", false},
		{"1 != 2", true},
		{"false != false", false},
		{"1 > 2", false},
		{"1 < 2", true},
		// TODO Compare integer and boolean
		//{"1 == false", false},
		//{"2 != true", false},
		{"!true", false},
		{"!!true", true},
		{"!1", false},
		{"!(if(false){5;})", true},
	}
	testRun(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []testCase{
		{"if(true){10}", 10},
		{"if(true){10}else{20}", 10},
		{"if(false){10}else{20}", 20},
		{"if((if(false){10})){10}else{20}", 20},
	}
	testRun(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []testCase{
		{"let one = 1;one", 1},
		{"let one = 1 let two = 2; one + two;", 3},
		{"let one = 1 let two = one +one; one + two;", 3},
	}
	testRun(t, tests)
}

func testRun(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		p := parse(tt.in)
		c := compiler.New()
		err := c.Compile(p)
		if err != nil {
			t.Fatalf("compiler got error: %s", err)
		}

		vm := New(c.ByteCode())
		err = vm.Run()

		if err != nil {
			t.Fatalf("vm.Run got error: %s", err)
		}

		stackElem := vm.LastPoppedStackElement()

		testExpectedObject(t, tt.in, tt.want, stackElem)
	}
}

func testRunWithError(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		p := parse(tt.in)
		c := compiler.New()
		err := c.Compile(p)
		if err != nil {
			t.Fatalf("compiler got error: %s", err)
		}

		vm := New(c.ByteCode())
		err = vm.Run()
		if err.Error() != tt.want.(error).Error() {
			t.Errorf("error is not %s, got %s", tt.want.(error).Error(), err.Error())
		}

	}
}

func testExpectedObject(
	t *testing.T,
	name string,
	want interface{},
	got object.Object,
) {
	t.Helper()
	switch want := want.(type) {
	case int:
		err := testIntegerObject(int64(want), got)
		if err != nil {
			t.Errorf("%s failed: %s", name, err)
		}
	case bool:
		err := testBooleanObject(bool(want), got)
		if err != nil {
			t.Errorf("%s failed: %s", name, err)
		}
	case *object.Null:
		if want != Null {
			t.Errorf("object is not Null: %T (%+v)", got, want)
		}
	}
}

func parse(in string) *ast.Program {
	l := lexer.New(in)
	p := parser.New(l)

	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	ret, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer.got=%T (%+v)", actual, actual)
	}
	if ret.Value != expected {
		return fmt.Errorf("object has wrong value. want=%d,got=%d", expected, ret.Value)
	}
	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	ret, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean.got=%T (%+v)", actual, actual)
	}
	if ret.Value != expected {
		return fmt.Errorf("object has wrong value. want=%v,got=%v", expected, ret.Value)
	}
	return nil
}

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
		{"1 == false", false},
		{"2 != true", true},
		{"!true", false},
		{"!!true", true},
		{"!1", false},
		{"!(if(false){5;})", true},
	}
	testRun(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []testCase{
		{`"monkey"`, "monkey"},
		{`"foo"+ "bar"`, "foobar"},
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

func TestArrayLiterals(t *testing.T) {
	tests := []testCase{
		{"[]", []int{}},
		{"[1,2,3]", []int{1, 2, 3}},
		{"[1+2,3*4]", []int{3, 12}},
	}
	testRun(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []testCase{
		{"{}", map[object.HashKey]int64{}},
	}
	testRun(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []testCase{
		{"[1,2,3][1]", 2},
		{"[[1,2,3]][0][0]", 1},
		{"[][0]", Null},
		{"[1][10]", Null},
		{"{1:1,2:2}[1]", 1},
		{"{1:1,2:2}[2]", 2},
		{"{1:1}[0]", Null},
		{"{}[0]", Null},
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

func TestCallingFunctions(t *testing.T) {
	tests := []testCase{
		{
			"let f = fn(){5 +10};f()",
			15},
		{
			`
				let one = fn(){1};
				let two = fn(){2};
				one() + two();
			`,
			3,
		},
		{
			`
				let one = fn(){1};
				let two = fn(){one()*2};
				fn(){two()+1}()
			`,
			3,
		},
		{
			`
				fn(){1;2}();
			`,
			2,
		},
		{
			`
				fn(){return 1;2}();
			`,
			1,
		},
	}
	testRun(t, tests)
}

func TestCallingFunctionsWithoutReturnValue(t *testing.T) {
	tests := []testCase{
		{
			`
				fn(){}();
			`,
			Null,
		},
		{
			`
				let f= fn(){};
				f();
			`,
			Null,
		},
		{
			`
				let a= fn(){};
				let b= fn(){a()};
				b();
			`,
			Null,
		},
	}
	testRun(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []testCase{
		{
			`
				fn(){let one = 1;one}();
			`,
			1,
		},
		{
			`
				let f = fn(){let one = 1;let two = 2;one + two;}
				f();
			`,
			3,
		},
		{
			`
				let f = fn(){let mon = "mon"; return mon + "key"}
				f();
			`,
			"monkey",
		},
		{
			`
				let f = fn(){let f = fn(){"monkey"};f}
				f()();
			`,
			"monkey",
		},
		{
			`
				let sum = fn(x,y){let z = x +y;z;}
				let outer = fn(){
					sum(1,2) + sum(3,4)		
				}
				outer();
			`,
			10,
		},
	}
	testRun(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []testCase{
		{
			`
				fn(two){let one = 1;one+two}(2);
			`,
			3,
		},
		{
			`
				let f = fn(three){let one = 1;let two = 2;one + two+three;}
				f(3);
			`,
			6,
		},
		{
			`
				let f = fn(key){let mon = "mon"; return mon + key}
				f("key");
			`,
			"monkey",
		},
	}
	testRun(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []testCase{
		{
			"fn(){1;}(1);",
			"wrong number of arguments: want=0, got=1",
		},
		{
			"fn(a){a;}();",
			"wrong number of arguments: want=1, got=0",
		},
		{
			"fn(x,y){x+y;}(1);",
			"wrong number of arguments: want=2, got=1",
		},
	}

	for _, tt := range tests {
		program := parse(tt.in)

		c := compiler.New()
		err := c.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(c.ByteCode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none")
		}

		if err.Error() != tt.want {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.want, err)
		}
	}
}

func TestClosures(t *testing.T) {
	tests := []testCase{
		{
			`
			let double = fn(x){
				fn(){
					return 2*x
				}
			}
			double(1)()`,
			2,
		},
		{
			`
			let head = fn(x){
				fn(){
					return first(x)
				}
			}
			head([10,1])()`,
			10,
		},
	}
	testRun(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []testCase{
		{
			`
			let f = fn(x){
				if (x < 2) { return x}
				return f(x-1) + f(x-2)
			}
			f(15)`,
			610,
		},
	}
	testRun(t, tests)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []testCase{
		{
			`len("")`, 0,
		},
		{
			`len([])`, 0,
		},
		{
			`len("four")`, 4,
		},
		{
			`len([1,2,3])`, 3,
		},
		{
			`len(1)`, &object.Error{Message: "argument to `len` not supported, got INTEGER"},
		},
		{
			`first([1,2,3])`, 1,
		},
		{
			`first([])`, Null,
		},
		{
			`first(1)`, &object.Error{Message: "argument to first must be ARRAY, got INTEGER"},
		},
		{
			`last([1,2,3])`, 3,
		},
		{
			`last([])`, Null,
		},
		{
			`last(1)`, &object.Error{Message: "argument to last must be ARRAY, got INTEGER"},
		},
		{
			`rest([1,2,3])`, []int{2, 3},
		},
		{
			`rest([])`, Null,
		},
		{
			`rest(1)`, &object.Error{Message: "argument to rest must be ARRAY, got INTEGER"},
		},
		{
			`puts("monkey")`, Null,
		},
		{
			`help()`, Null,
		},
		{
			`exit()`, Null,
		},
	}
	testRun(t, tests)
}

func TestIssue001(t *testing.T) {
	tests := []testCase{
		{"return 1;", 1},
		{"if(true){return \"x\";}", "x"},
		{"fn(){if(true){return true;};}();", true},
	}
	testRun(t, tests)
}

func testRun(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
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
		})
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
	case string:
		err := testStringObject(string(want), got)
		if err != nil {
			t.Errorf("%s failed: %s", name, err)
		}
	case []int:
		array, ok := got.(*object.Array)
		if !ok {
			t.Errorf("object not Array:%T (%+v)", got, got)
		}
		if len(array.Elements) != len(want) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(want), len(array.Elements))
		}
		for key, value := range want {
			err := testIntegerObject(int64(value), array.Elements[key])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := got.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", got, got)
		}
		if len(hash.Pairs) != len(want) {
			t.Errorf("wrong num of pairs. want=%d, got=%d", len(want), len(hash.Pairs))
		}
		for key, value := range want {
			pair, ok := hash.Pairs[key]
			if !ok {
				t.Errorf("no pair for given key `%v` in pairs", key.Value)
			}
			err := testIntegerObject(value, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case *object.Null:
		if got != Null {
			t.Errorf("object is not Null: %T (%+v)", got, want)
		}
	case *object.Error:
		{
			err, ok := got.(*object.Error)
			if !ok {
				t.Errorf("object is not Error: %T (%+v)", got, want)
			}
			if err.Message != want.Message {
				t.Errorf("wrong error message. want=%q, got=%q", want.Message, err.Message)
			}
		}
	default:
		t.Errorf("test is not implemented. %T (%+v)", got, want)
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

func testStringObject(expected string, actual object.Object) error {
	ret, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not Boolean.got=%T (%+v)", actual, actual)
	}
	if ret.Value != expected {
		return fmt.Errorf("object has wrong value. want=%v,got=%v", expected, ret.Value)
	}
	return nil
}

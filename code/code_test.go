package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       OperandCode
		operands []int
		expected []byte
	}{
		{Constant, []int{1}, []byte{byte(Constant), 0, 1}},
		{Constant, []int{65534}, []byte{byte(Constant), 255, 254}},
		{Add, []int{}, []byte{byte(Add)}},
		{GetLocal, []int{255}, []byte{byte(GetLocal), 255}},
		{OperandCode(255), nil, []byte{}},
		{Closure, []int{65534, 255}, []byte{byte(Closure), 255, 254, 255}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d", len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(Add),
		Make(Constant, 2),
		Make(Constant, 65535),
		Make(Pop),
		Make(Sub),
		Make(Mul),
		Make(Div),
		Make(Equal),
		Make(NotEqual),
		Make(GreaterThan),
		Make(True),
		Make(False),
		Make(Minus),
		Make(Bang),
		Make(Jump, 1),
		Make(JumpNotTruthy, 2),
		Make(GetLocal, 1),
		Make(Null),
		Make(Closure, 65535, 255),
	}
	expected := `0000 Add
0001 Constant 2
0004 Constant 65535
0007 Pop
0008 Sub
0009 Mul
0010 Div
0011 Equal
0012 NotEqual
0013 GreaterThan
0014 True
0015 False
0016 Minus
0017 Bang
0018 Jump 1
0021 JumpNotTruthy 2
0024 GetLocal 1
0026 Null
0027 Closure 65535 255
`

	concatted := Instructions{}

	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot= %q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        OperandCode
		operands  []int
		bytesRead int
	}{
		{Constant, []int{65535}, 2},
		{Add, []int{}, 0},
		{Sub, []int{}, 0},
		{Mul, []int{}, 0},
		{Div, []int{}, 0},
		{Pop, []int{}, 0},
		{True, []int{}, 0},
		{Equal, []int{}, 0},
		{NotEqual, []int{}, 0},
		{GreaterThan, []int{}, 0},
		{Minus, []int{}, 0},
		{Bang, []int{}, 0},
		{Jump, []int{}, 2},
		{JumpNotTruthy, []int{}, 2},
		{Closure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := LookUp(byte(tt.op))

		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])

		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d,got=%d", want, operandsRead[i])
			}
		}
	}
}

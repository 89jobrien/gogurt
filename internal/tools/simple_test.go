package tools

import "testing"

func TestUppercaseTool(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:  "Simple lowercase string",
			input: `{"Text":"gogurt"}`,
			want:  "GOGURT",
		},
		{
			name:  "Mixed case string",
			input: `{"Text":"GoGuRt"}`,
			want:  "GOGURT",
		},
		{
			name:  "Empty string",
			input: `{"Text":""}`,
			want:  "",
		},
		{
			name:    "Invalid JSON",
			input:   `{"Text": "invalid`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UppercaseTool.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UppercaseTool.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("UppercaseTool.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcatenateTool(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:  "Two simple strings",
			input: `{"a":"hello", "b":" world"}`,
			want:  "hello world",
		},
		{
			name:  "First string is empty",
			input: `{"a":"", "b":"world"}`,
			want:  "world",
		},
		{
			name:  "Second string is empty",
			input: `{"a":"hello", "b":""}`,
			want:  "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConcatenateTool.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConcatenateTool.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ConcatenateTool.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name    string
		args    ReverseArgs
		want    string
		wantErr bool
	}{
		{
			name: "Simple string",
			args: ReverseArgs{Text: "hello"},
			want: "olleh",
		},
		{
			name: "Empty string",
			args: ReverseArgs{Text: ""},
			want: "",
		},
		{
			name: "Palindrome",
			args: ReverseArgs{Text: "racecar"},
			want: "racecar",
		},
		{
			name: "String with spaces",
			args: ReverseArgs{Text: "hello world"},
			want: "dlrow olleh",
		},
		{
			name: "String with unicode characters",
			args: ReverseArgs{Text: "hello, 世界"},
			want: "界世 ,olleh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Reverse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseTool(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:  "Simple reverse",
			input: `{"Text":"gogurt"}`,
			want:  "trugog",
		},
		{
			name:  "Empty string",
			input: `{"Text":""}`,
			want:  "",
		},
		{
			name:    "Invalid JSON",
			input:   `{"Text":}`,
			wantErr: true,
		},
		{
			name:  "Unicode string",
			input: `{"Text":"你好"}`,
			want:  "好你",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReverseTool.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReverseTool.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReverseTool.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPalindromeTool(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{`{"Text":"racecar"}`, true},
		{`{"Text":"abc"}`, false},
		{`{"Text":""}`, true},
		{`{"Text":"A man a plan a canal Panama"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			res, err := PalindromeTool.Call(tt.input)
			if err != nil {
				t.Fatalf("error on input %s: %v", tt.input, err)
			}
			if res != tt.expect {
				t.Errorf("for input %s, expected %v got %v", tt.input, tt.expect, res)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		args NumInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"both positive", args{NumInput{1, 2}}, 3, false},
		{"one zero", args{NumInput{0, 8}}, 8, false},
		{"both zero", args{NumInput{0, 0}}, 0, false},
		{"negative and positive", args{NumInput{-5, 3}}, -2, false},
		{"both negative", args{NumInput{-5, -10}}, -15, false},
		{"large numbers", args{NumInput{1000000, 2000000}}, 3000000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(NumInput(tt.args.args))
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	type args struct {
		args NumInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"positive", args{NumInput{5, 2}}, 3, false},
		{"zero", args{NumInput{5, 0}}, 5, false},
		{"reverse", args{NumInput{2, 5}}, -3, false},
		{"negative and positive", args{NumInput{-7, 3}}, -10, false},
		{"both negative", args{NumInput{-7, -3}}, -4, false},
		{"large numbers", args{NumInput{1000000, 500000}}, 500000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subtract(NumInput(tt.args.args))
			if (err != nil) != tt.wantErr {
				t.Errorf("Subtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	type args struct {
		args NumInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"positive", args{NumInput{2, 3}}, 6, false},
		{"zero", args{NumInput{8, 0}}, 0, false},
		{"negative and positive", args{NumInput{-2, 3}}, -6, false},
		{"both negative", args{NumInput{-2, -3}}, 6, false},
		{"large numbers", args{NumInput{1000, 2000}}, 2000000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Multiply(NumInput(tt.args.args))
			if (err != nil) != tt.wantErr {
				t.Errorf("Multiply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	type args struct {
		args NumInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"positive", args{NumInput{6, 3}}, 2, false},
		{"zero numerator", args{NumInput{0, 3}}, 0, false},
		{"negative numerator", args{NumInput{-6, 3}}, -2, false},
		{"both negative", args{NumInput{-6, -3}}, 2, false},
		{"div by zero", args{NumInput{5, 0}}, 0, true},
		{"zero divided by zero", args{NumInput{0, 0}}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(NumInput(tt.args.args))
			if (err != nil) != tt.wantErr {
				t.Errorf("Divide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}

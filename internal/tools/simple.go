package tools

import (
	"errors"
	"reflect"
	"strings"
)

// --- Uppercase Tool ---

type UppercaseArgs struct {
	Text string `json:"Text"`
}

func Uppercase(args UppercaseArgs) (string, error) {
	return strings.ToUpper(args.Text), nil
}

var UppercaseTool = &Tool{
	Name:        "uppercase",
	Description: "Converts a string to uppercase.",
	Func:        reflect.ValueOf(Uppercase),
	InputSchema: GenInputSchema(reflect.TypeOf(UppercaseArgs{})),
	Example:     `{"Text":"hello world"}`,
	Metadata:    map[string]any{"author": "joe", "name": "Uppercase", "category": "text", "version": "0.1", "deprecated": false},
}

// --- Concatenate Tool ---

type ConcatArgs struct {
	A string `json:"a"`
	B string `json:"b"`
}

func Concatenate(args ConcatArgs) (string, error) {
	if args.A == "" {
		return args.B, nil
	}
	if args.B == "" {
		return args.A, nil
	}
	return args.A + args.B, nil
}

var ConcatenateTool = &Tool{
	Name:        "concatenate",
	Description: "Joins two strings together.",
	Func:        reflect.ValueOf(Concatenate),
	InputSchema: GenInputSchema(reflect.TypeOf(ConcatArgs{})),
	Example:     `{"a":"hello", "b":" world"}`,
	Metadata:    map[string]any{"author": "joe", "category": "text", "version": "0.1", "deprecated": false},
}

// --- Reverse Tool ---

type ReverseArgs struct {
	Text string `json:"Text"`
}

func Reverse(args ReverseArgs) (string, error) {
	r := []rune(args.Text)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

var ReverseTool = &Tool{
	Name:        "reverse",
	Description: "Reverses a string.",
	Func:        reflect.ValueOf(Reverse),
	InputSchema: GenInputSchema(reflect.TypeOf(ReverseArgs{})),
	Example:     `{"Text":"foo"}`,
	Metadata:    map[string]any{"author": "joe", "category": "text", "version": "0.1", "deprecated": false},
}

// --- Palindrome Tool ---

type PalindromeArgs struct {
	Text string `json:"Text"`
}

func Palindrome(args PalindromeArgs) (bool, error) {
	r := []rune(args.Text)
	n := len(r)
	for i := 0; i < n/2; i++ {
		if r[i] != r[n-1-i] {
			return false, nil
		}
	}
	return true, nil
}

var PalindromeTool = &Tool{
	Name:        "palindrome",
	Description: "Checks if the given string is a palindrome (reads the same forwards and backwards).",
	Func:        reflect.ValueOf(Palindrome),
	InputSchema: GenInputSchema(reflect.TypeOf(PalindromeArgs{})),
	Example:     `{"Text":"racecar"}`,
	Metadata:    map[string]any{"author": "joe", "category": "text", "version": "0.1", "deprecated": false},
}

// --- Calculation Tools ---

// NumInput is used as input for arithmetic operations.
type NumInput struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Add returns the sum of a and b.
func Add(args NumInput) (int, error) {
	return args.A + args.B, nil
}

var AddTool = &Tool{
	Name:        "add",
	Description: "Returns the sum of a and b.",
	Func:        reflect.ValueOf(Add),
	InputSchema: GenInputSchema(reflect.TypeOf(NumInput{})),
	Example:     `{"a":3, "b":4}`,
	Metadata:    map[string]any{"author": "joe", "category": "math", "version": "0.1", "deprecated": false},
}

// Subtract returns the difference of a and b.
func Subtract(args NumInput) (int, error) {
	return args.A - args.B, nil
}

var SubtractTool = &Tool{
	Name:        "subtract",
	Description: "Returns the difference of a and b.",
	Func:        reflect.ValueOf(Subtract),
	InputSchema: GenInputSchema(reflect.TypeOf(NumInput{})),
	Example:     `{"a":7, "b":2}`,
	Metadata:    map[string]any{"author": "joe", "category": "math", "version": "0.1", "deprecated": false},
}

// Multiply returns the product of a and b.
func Multiply(args NumInput) (int, error) {
	return args.A * args.B, nil
}

var MultiplyTool = &Tool{
	Name:        "multiply",
	Description: "Returns the product of a and b.",
	Func:        reflect.ValueOf(Multiply),
	InputSchema: GenInputSchema(reflect.TypeOf(NumInput{})),
	Example:     `{"a":3, "b":5}`,
	Metadata:    map[string]any{"author": "joe", "category": "math", "version": "0.1", "deprecated": false},
}

// Divide returns the integer division of a by b.
func Divide(args NumInput) (int, error) {
	if args.B == 0 {
		return 0, errors.New("division by zero")
	}
	return args.A / args.B, nil
}

var DivideTool = &Tool{
	Name:        "divide",
	Description: "Returns the integer division of a by b.",
	Func:        reflect.ValueOf(Divide),
	InputSchema: GenInputSchema(reflect.TypeOf(NumInput{})),
	Example:     `{"a":14, "b":2}`,
	Metadata:    map[string]any{"author": "joe", "category": "math", "version": "0.1", "deprecated": false},
}

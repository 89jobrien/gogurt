package tools

// import (
// 	"gogurt/tools" // update this import for your codebase
// )

type AddInput struct {
    A int `json:"a"`
    B int `json:"b"`
}

func Add(args AddInput) (int, error) {
    return args.A + args.B, nil
}

func Subtract(args AddInput) (int, error) {
    return args.A - args.B, nil
}

func Multiply(args AddInput) (int, error) {
	return args.A * args.B, nil
}

func Divide(args AddInput) (int, error) {
	return args.A / args.B, nil
}


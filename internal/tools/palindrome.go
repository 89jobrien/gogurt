package tools

import "reflect"

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
    InputSchema: map[string]any{
        "type": "object",
        "properties": map[string]any{
            "Text": map[string]any{"type": "string"},
        },
        "required": []string{"Text"},
    },
    Example:  `{"Text":"racecar"}`,
    Metadata: map[string]any{"category": "text"},
}
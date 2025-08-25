package tools

import "testing"

func TestPalindromeTool(t *testing.T) {
    tests := []struct {
        input  string
        expect bool
    }{
        {`{"Text":"racecar"}`, true},
        {`{"Text":"abc"}`, false},
        {`{"Text":""}`, true},
        {`{"Text":"A man a plan a canal Panama"}`, false}, // spacing/cases matter in this basic version
    }
    for _, tc := range tests {
        res, err := PalindromeTool.Call(tc.input)
        if err != nil {
            t.Fatalf("error on input %s: %v", tc.input, err)
        }
        if res != tc.expect {
            t.Errorf("for input %s, expected %v got %v", tc.input, tc.expect, res)
        }
    }
}
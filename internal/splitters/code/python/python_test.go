package python

import (
	"reflect"
	"testing"
)

func TestSplit_NoFunctionsOrClasses(t *testing.T) {
	content := `print("hello")`

	got := Split(content)
	want := []string{`print("hello")`}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %v, want %v", got, want)
	}
}

func TestSplit_SingleFunction(t *testing.T) {
	content := `def greet():
    print("hi")
`

	got := Split(content)
	want := []string{
		`def greet():
    print("hi")
`,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %v, want %v", got, want)
	}
}

func TestSplit_MultipleFunctionsAndClasses(t *testing.T) {
	content := `def foo():
    pass

class Bar:
    def baz(self):
        pass
`

	got := Split(content)
	want := []string{
		`def foo():
    pass

`,
		`class Bar:
    def baz(self):
        pass
`,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %v, want %v", got, want)
	}
}

func TestSplit_TextBeforeFirstDef(t *testing.T) {
	content := `# Some module comments
import os

def foo():
    pass

def bar():
    pass
`

	got := Split(content)
	want := []string{
		`# Some module comments
import os

`,
		`def foo():
    pass

`,
		`def bar():
    pass
`,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %v, want %v", got, want)
	}
}

func TestSplit_FunctionWithoutTrailingNewline(t *testing.T) {
	content := `def foo():
    return 1`

	got := Split(content)
	want := []string{
		`def foo():
    return 1`,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %v, want %v", got, want)
	}
}

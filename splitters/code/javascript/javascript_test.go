package javascript

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestSplit_NoFunctionsOrClasses(t *testing.T) {
	content := `console.log("hello world");`

	got := Split(content)
	want := []string{content}
	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

func TestSplit_SimpleFunction(t *testing.T) {
	content := `function greet() {
  console.log("hi");
}
`
	got := Split(content)
	want := []string{content}
	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

func TestSplit_ClassDeclaration(t *testing.T) {
	content := `class Person {
  constructor(name) {
    this.name = name;
  }
}
`
	got := Split(content)
	want := []string{content}
	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

func TestSplit_WithExportAndAsync(t *testing.T) {
	content := `export async function load() {
  return fetch("/api");
}

export default class Service {}
`
	got := Split(content)
	want := []string{
		"export async function load() {\n  return fetch(\"/api\");\n}\n",
		"\nexport default class Service {}\n",
	}
	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		fmt.Printf("\nEXPECTED:\n%q\n", want[0])
		fmt.Printf("GOT:\n%q\n\n", got)
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

func TestSplit_WithConstLetVar(t *testing.T) {
	content := `const PI = 3.14;

let counter = 0;

var legacy = true;
`
	got := Split(content)
	want := []string{
		"const PI = 3.14;\n",
		"\nlet counter = 0;\n",
		"\nvar legacy = true;\n",
	}

	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		fmt.Printf("\nEXPECTED:\n%q\n", want[0])
		fmt.Printf("GOT:\n%q\n\n", got)
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

func TestSplit_ImportsAndPreludeText(t *testing.T) {
	content := `// utilities
import fs from "fs";

export function util() {
  return fs.readFileSync("file.txt");
}
`
	got := Split(content)
	want := []string{
		"// utilities\nimport fs from \"fs\";\n",
		"\nexport function util() {\n  return fs.readFileSync(\"file.txt\");\n}\n",
	}

	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		fmt.Printf("\nEXPECTED:\n%q\n", want[0])
		fmt.Printf("GOT:\n%q\n\n", got)
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

func TestSplit_NoTrailingNewline(t *testing.T) {
	content := `export function lastOne() {
  return true;
}`
	got := Split(content)
	want := []string{content}
	if strings.Join(got, "") != content {
		t.Fatalf("joined chunks != original\njoined:\n%q\noriginal:\n%q", strings.Join(got, ""), content)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split() = %+v, want %+v", got, want)
	}
}

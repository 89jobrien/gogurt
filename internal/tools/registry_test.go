package tools

import (
	"reflect"
	"strings"
	"testing"
)

type dummyArgs struct {
	Field string `json:"Field"`
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if len(r.tools) != 0 {
		t.Errorf("expected empty registry, got %d", len(r.tools))
	}
}

func makeTestTool(name, category string) *Tool {
	tool := &Tool{
		Name:        name,
		Description: "desc: " + name,
		Func:        reflect.ValueOf(func(_ dummyArgs) (string, error) { return "ok", nil }),
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{"Field": map[string]any{"type": "string"}}, "required": []string{"Field"}},
		Example:     `{"Field":"val"}`,
		Metadata:    map[string]any{"category": category},
	}
	return tool
}

func TestRegisterToolNil(t *testing.T) {
	r := NewRegistry()
	err := r.Register(nil)
	if err == nil || err.Error() != "cannot register nil tool" {
		t.Errorf("Register(nil) error = %v, expected 'cannot register nil tool'", err)
	}
}

func TestRegisterToolInvalidName(t *testing.T) {
	r := NewRegistry()
	tool := makeTestTool("   ", "cat")
	err := r.Register(tool)
	if err == nil || !strings.Contains(err.Error(), "tool name invalid") {
		t.Errorf("Register invalid name: error = %v, expected invalid name error", err)
	}
	tool = makeTestTool("invalid name", "cat")
	err = r.Register(tool)
	if err == nil || !strings.Contains(err.Error(), "tool name invalid") {
		t.Errorf("Register invalid name with space: error = %v, expected invalid name error", err)
	}
}

func TestRegisterToolDup(t *testing.T) {
	r := NewRegistry()
	t1 := makeTestTool("foo", "cat")
	if err := r.Register(t1); err != nil {
		t.Fatalf("Register foo: %v", err)
	}
	t2 := makeTestTool("foo", "cat")
	err := r.Register(t2)
	if err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Errorf("Register duplicate: error = %v, expected already registered error", err)
	}
}

func TestRegisterBatch(t *testing.T) {
	r := NewRegistry()
	t1 := makeTestTool("a", "alpha")
	t2 := makeTestTool("b", "beta")
	t3 := makeTestTool("a", "alpha") // duplicate name
	errs := r.RegisterBatch([]*Tool{t1, t2, t3})
	if len(errs) != 1 {
		t.Errorf("expected one error for duplicate registration, got %d", len(errs))
	}
}

func TestGet(t *testing.T) {
	r := NewRegistry()
	t1 := makeTestTool("one", "x")
	r.Register(t1)
	if r.Get("one") == nil {
		t.Errorf("Get(one) should return tool")
	}
	if r.Get("noname") != nil {
		t.Errorf("Get(noname) should return nil")
	}
	if r.Get("a bad name") != nil {
		t.Errorf("Get('a bad name') should return nil for invalid name")
	}
}

func TestCall(t *testing.T) {
	r := NewRegistry()
	tool := &Tool{
		Name:        "t1",
		Description: "desc: t1",
		Func:        reflect.ValueOf(func(args dummyArgs) (string, error) { return "called:t1:" + args.Field, nil }),
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{"Field": map[string]any{"type": "string"}}, "required": []string{"Field"}},
		Example:     `{"Field":"abc"}`,
		Metadata:    map[string]any{"category": "cat"},
	}
	r.Register(tool)

	val, err := r.Call("t1", `{"Field":"abc"}`)
	if err != nil {
		t.Errorf("Call(t1) error = %v", err)
	}
	if val != "called:t1:abc" {
		t.Errorf("Call(t1) value = %v, want called:t1:abc", val)
	}
	val2, err2 := r.Call("notfound", `{"Field":"test"}`)
	if err2 == nil || !strings.Contains(err2.Error(), "not found") {
		t.Errorf("Call(notfound) error = %v, want not found error", err2)
	}
	if val2 != nil {
		t.Errorf("Call(notfound) value = %v, want nil", val2)
	}
}

func TestList(t *testing.T) {
	r := NewRegistry()
	names := r.List()
	if len(names) != 0 {
		t.Errorf("List on empty registry should be empty, got %v", names)
	}
	r.Register(makeTestTool("a", "cat"))
	r.Register(makeTestTool("b", "cat"))
	got := r.List()
	want := []string{"a", "b"}
	if !containsAllStrings(got, want) {
		t.Errorf("List got %v, want %v", got, want)
	}
}

func containsAllStrings(got, want []string) bool {
	for _, w := range want {
		found := false
		for _, g := range got {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestListTools(t *testing.T) {
	r := NewRegistry()
	t1 := makeTestTool("toolA", "alpha")
	t2 := makeTestTool("toolB", "beta")
	r.RegisterBatch([]*Tool{t1, t2})
	tools := r.ListTools()
	names := []string{tools[0].Name, tools[1].Name}
	if !containsAllStrings(names, []string{"toolA", "toolB"}) {
		t.Errorf("ListTools names = %v, expected toolA, toolB", names)
	}
}

func TestGetByCategory(t *testing.T) {
	r := NewRegistry()
	r.Register(makeTestTool("a", "cat"))
	r.Register(makeTestTool("b", "dog"))
	r.Register(makeTestTool("c", "cat"))
	got := r.GetByCategory("cat")
	if len(got) != 2 || !containsToolByName(got, "a") || !containsToolByName(got, "c") {
		t.Errorf("GetByCategory(cat) = %v", toolListNames(got))
	}
	got2 := r.GetByCategory("dog")
	if len(got2) != 1 || got2[0].Name != "b" {
		t.Errorf("GetByCategory(dog) = %v", toolListNames(got2))
	}
	got3 := r.GetByCategory("fish")
	if len(got3) != 0 {
		t.Errorf("GetByCategory(fish) should be empty, got %v", toolListNames(got3))
	}
}

func containsToolByName(ts []*Tool, name string) bool {
	for _, t := range ts {
		if t.Name == name {
			return true
		}
	}
	return false
}

func toolListNames(ts []*Tool) []string {
	names := make([]string, 0, len(ts))
	for _, t := range ts {
		names = append(names, t.Name)
	}
	return names
}

func TestStats(t *testing.T) {
	r := NewRegistry()
	r.Register(makeTestTool("a", "cat"))
	r.Register(makeTestTool("b", "dog"))
	r.Register(makeTestTool("c", "dog"))
	st := r.Stats()
	if st.Count != 3 {
		t.Errorf("Stats.Count = %d, want 3", st.Count)
	}
	if !containsAllStrings(st.ToolNames, []string{"a", "b", "c"}) {
		t.Errorf("Stats.ToolNames = %v", st.ToolNames)
	}
	if !containsAllStrings(st.Categories, []string{"cat", "dog"}) {
		t.Errorf("Stats.Categories = %v", st.Categories)
	}
	if st.HasDups {
		t.Errorf("Stats.HasDups: expected false")
	}
	for _, cat := range []string{"cat", "dog"} {
		if !st.HasCategory[cat] {
			t.Errorf("Stats.HasCategory[%q] not found", cat)
		}
	}
}

func TestRegistry_ZeroStates(t *testing.T) {
	r := NewRegistry()
	if got := r.List(); len(got) != 0 {
		t.Errorf("empty List: %v", got)
	}
	if got := r.ListTools(); len(got) != 0 {
		t.Errorf("empty ListTools: %v", got)
	}
	if got := r.GetByCategory("foo"); len(got) != 0 {
		t.Errorf("empty GetByCategory: %v", got)
	}
	st := r.Stats()
	if st.Count != 0 {
		t.Errorf("Stats.Count = %d, want 0", st.Count)
	}
	if len(st.ToolNames) != 0 {
		t.Errorf("Stats.ToolNames = %v, want 0", st.ToolNames)
	}
	if len(st.Categories) != 0 {
		t.Errorf("Stats.Categories = %v, want 0", st.Categories)
	}
	if st.HasDups {
		t.Errorf("Stats.HasDups: expected false")
	}
}

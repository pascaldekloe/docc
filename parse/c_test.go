package parse

import (
	"reflect"
	"testing"
)

//go:generate ragel -Z c.rl

type goldenParse struct {
	lines    int
	feed     string
	comments []string
}

var goldenCs = []goldenParse{
	{3, "// line comment\nint i, j;\n", []string{"// line comment\n"}},
	{3, "\npre// trailing comment\n", []string{"// trailing comment\n"}},
	{2, "pre/*inner comment*/fix\n", []string{"/*inner comment*/"}},
	{3, "/* block\ncomment */\n", []string{"/* block\ncomment */"}},
	{4, "int main(void) {\n//inner comment line\n}\n", nil},
	{1, "int main(void) {a /*inner comment block*/ ignore}", nil},
	{1, `"string"`, nil},
	{1, `"a\"b"`, nil},
	{1, `"//not a comment"`, nil},
	{1, `'c'`, nil},
	{1, `'\''`, nil},
	{1, `'{'`, nil},
}

func TestC(t *testing.T) {
	for _, gold := range goldenCs {
		lineCount, comments := C([]byte(gold.feed))
		if lineCount != gold.lines || !reflect.DeepEqual(comments, gold.comments) {
			t.Errorf("%q: got (%d, %q), want (%d, %q)", gold.feed, lineCount, comments, gold.lines, gold.comments)
		}
	}
}

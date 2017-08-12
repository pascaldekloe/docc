package parse // import "docc.io/source/parse"

import (
	"io/ioutil"
	"strings"
	"testing"

	"docc.io/source"
)

var defTypes = []string{
	`extern int var_decl;`,
	`int var_def;`,
	`int func_decl(void);`,
	`int func_def(void) {`,
}

var commentTypes = []string{
	`// single line`,
	`// line 1
// line 2`,
	`/* block */`,
	`/* block 1
block 2 */`,
}

func TestComments(t *testing.T) {
	for _, def := range defTypes {
		for _, comment := range commentTypes {
			feed := comment + "\n" + def

			c := make(chan *source.Decl)
			go C([]byte(feed), c)

			var got *source.Decl
			for d := range c {
				if got == nil {
					got = d
				} else {
					t.Errorf("%q: got redundant result %#v", feed, d)
				}
			}
			if got == nil {
				t.Errorf("%q: no results", feed)
				continue
			}

			if want := strings.Count(comment, "\n") + 2; got.LineNo != want {
				t.Errorf("%q: got line number %d, want %d", feed, got.LineNo, want)
			}
			if got.Source != def {
				t.Errorf("%q: got source %q, want %q", feed, got.Source, def)
			}
			if got.Comment != comment {
				t.Errorf("%q: got comment %q, want %q", feed, got.Comment, comment)
			}
		}
	}
}

func TestStdioH(t *testing.T) {
	data, err := ioutil.ReadFile("/usr/include/stdio.h")
	if err != nil {
		t.Fatal(err)
	}

	c := make(chan *source.Decl, 10)
	go C(data, c)
	for def := range c {
		t.Logf("stdio.h:%d: %q\n\t%q\n\n", def.LineNo, def.Source, def.Comment)
	}
}

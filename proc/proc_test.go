package proc // import "docc.io/source/proc"

import "testing"

func TestRepo1(t *testing.T) {
	files, err := Repo("./testdata/repo1")
	if err != nil {
		t.Fatal(err)
	}

	wantFiles := map[string]bool{
		"x.h":     true,
		"y.h":     true,
		"sub/z.h": true,
	}
	wantVars := map[string]bool{
		"extern int x;":          true,
		"extern const char z[];": true,
	}
	wantFuncs := map[string]bool{
		"int y(void);": true,
	}

	for f := range files {
		switch want, ok := wantFiles[f.Path]; {
		case !ok:
			t.Errorf("got unknown file %q", f.Path)
		case !want:
			t.Errorf("got redundant file %q", f.Path)
		default:
			wantFiles[f.Path] = false
		}

		for _, decl := range f.Vars {
			switch want, ok := wantVars[decl.Source]; {
			case !ok:
				t.Errorf("got unknown variable %q", decl.Source)
			case !want:
				t.Errorf("got redundant variable %q", decl.Source)
			default:
				wantVars[decl.Source] = false
			}
		}

		for _, decl := range f.Funcs {
			switch want, ok := wantFuncs[decl.Source]; {
			case !ok:
				t.Errorf("got unknown function %q", decl.Source)
			case !want:
				t.Errorf("got redundant function %q", decl.Source)
			default:
				wantFuncs[decl.Source] = false
			}
		}

	}

	for path, want := range wantFiles {
		if want {
			t.Errorf("file %q not seen", path)
		}
	}
	for s, want := range wantVars {
		if want {
			t.Errorf("variable %q not seen", s)
		}
	}
	for s, want := range wantFuncs {
		if want {
			t.Errorf("function %q not seen", s)
		}
	}
}

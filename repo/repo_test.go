package repo // import "docc.io/source/repo"

import "testing"

func TestValidNameSeg(t *testing.T) {
	var valid = []string{
		"a",
		"β",
		"1",
		"Ⅱ",
		"_.",
		"Z-",
	}
	for _, s := range valid {
		if err := ValidNameSeg(s); err != nil {
			t.Errorf("%q: %s", s, err)
		}
	}

	var invalid = []string{
		"",
		".",
		"..",
		".a",
		"-",
		"-a",
		"\x00",
		"a\x00z",
		"\xff",
		"a\xffz",
		"🌲",
	}
	for _, s := range invalid {
		if err := ValidNameSeg(s); err == nil {
			t.Errorf("no error for %q", s)
		}
	}
}

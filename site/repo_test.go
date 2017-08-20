package site // import "docc.io/source/site"

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"docc.io/source"
	"docc.io/source/repo"
)

func init() {
	// mock repository lookup
	latestFiles = func(id repo.Name) []*source.File {
		return []*source.File{
			{
				Vars: []*source.Decl{
					{LineNo: 1, Source: "extern int g", Comment: "// !"},
				},
				Funcs: []*source.Decl{
					{LineNo: 42, Source: "int f(void)", Comment: "// ?"},
				},
			},
		}
	}
}

func TestGitHubDedupe(t *testing.T) {
	srv := httptest.NewServer(Router)
	defer srv.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(srv.URL + "/github.com/UPPER/CASE")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 301 {
		t.Errorf("got status code %d, want 301", resp.StatusCode)
	}
	if got, want := resp.Header.Get("Location"), "/github.com/upper/case"; got != want {
		t.Errorf("got location %q, want %q", got, want)
	}
}

func TestGitHubMisses(t *testing.T) {
	srv := httptest.NewServer(Router)
	defer srv.Close()

	paths := []string{
		"/one",
		"/one/two/tree",
		"/hidden/-arg",
		"/.hidden/a",
	}

	for _, p := range paths {
		resp, err := http.Get(srv.URL + "/github.com" + p)
		if err != nil {
			t.Errorf("%q: %s", p, err)
			continue
		}

		if resp.StatusCode != 404 {
			t.Errorf("%q: got status code %d, want 404", p, resp.StatusCode)
		}
	}
}

package site // import "docc.io/source/site"

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"

	"docc.io/source/repo"
)

type Display struct {
	ID repo.Name
	// Version control location.
	URL string

	RepoLabel string
	// URL to the main page.
	RepoLink string

	AccountLabel string
	// URL to the profile page.
	AccountLink string

	HostLabel string
	// URL to the main page.
	HostLink string
}

const GitHubLabel = "GitHub"
const GitHubLink = "https://github.com/"

func gitHubGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	accountParam, repoParam := params[0].Value, params[1].Value

	if err := repo.ValidNameSeg(accountParam); err != nil {
		notFound(w, fmt.Sprintf("Rejecting %q: %s", accountParam, err))
		return
	}
	if err := repo.ValidNameSeg(repoParam); err != nil {
		notFound(w, fmt.Sprintf("Rejecting %q: %s", repoParam, err))
		return
	}

	// prevent duplicates with case insensitivity
	if lr, la := strings.ToLower(repoParam), strings.ToLower(accountParam); lr != repoParam || la != accountParam {
		http.Redirect(w, r, fmt.Sprintf("/github.com/%s/%s", url.PathEscape(la), url.PathEscape(lr)), http.StatusMovedPermanently)
		return
	}

	// clear

	url := GitHubLink + accountParam + "/" + repoParam + ".git"
	d := &Display{
		URL:          url,
		AccountLabel: accountParam,
		RepoLabel:    repoParam,
		HostLabel:    GitHubLabel,
		HostLink:     GitHubLink,
	}
	d.RepoLink = url[:len(url)-4]
	d.AccountLink = url[:len(url)-5-len(repoParam)]
	d.ID = repo.Name(d.RepoLink[len(GitHubLink):])

	h := w.Header()
	h.Set("Content-Type", "text/html;charset=utf-8")
	h.Set("Content-Language", "en")

	repoPage.Execute(w, d)
}

var repoPage = template.Must(template.New("repo-html").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html;charset=utf-8">
<title>docc.io: {{.ID}}</title>
</head>
<body>

<h1><a href="{{.RepoLink}}">{{.RepoLabel}}</a> from <a href="{{.AccountLink}}">{{.AccountLabel}}</a>@<a href="{{.HostLink}}">{{.HostLabel}}</a></h>

</body>
`))

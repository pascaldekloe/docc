package site // import "docc.io/source/site"

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"

	"docc.io/source"
	"docc.io/source/proc"
	"docc.io/source/repo"
)

var latestFiles = proc.Latest

// Display holds repository information.
type Display struct {
	ID repo.Name

	RepoLabel string
	// URL to the main page.
	RepoLink string

	AccountLabel string
	// URL to the profile page.
	AccountLink string

	HostLabel string
	// URL to the main page.
	HostLink string

	Files []*source.File
}

const gitHubLabel = "GitHub"
const gitHubLink = "https://github.com/"

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

	d := &Display{
		AccountLabel: accountParam,
		RepoLabel:    repoParam,
		RepoLink:     gitHubLink + accountParam + "/" + repoParam,
		HostLabel:    gitHubLabel,
		HostLink:     gitHubLink,
	}
	d.AccountLink = d.RepoLink[:len(d.RepoLink)-len(repoParam)-1]
	d.ID = repo.Name(d.RepoLink[8:])

	d.Files = latestFiles(d.ID)

	h := w.Header()
	h.Set("Content-Type", "text/html;charset=utf-8")
	h.Set("Content-Language", "en")

	if err := repoPage.Execute(w, d); err != nil {
		log.Printf("%d: template render: %s", d.ID, err)
	}
}

var repoPage = template.Must(template.New("repo-html").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html;charset=utf-8">
<title>docc.io: {{.ID}}</title>
</head>
<body>

<h1><a href="{{.RepoLink}}">{{.RepoLabel}}</a> from <a href="{{.AccountLink}}">{{.AccountLabel}}</a>@<a href="{{.HostLink}}">{{.HostLabel}}</a></h>

{{range $index, $file := .Files}}

{{range $indexx, $var := .Vars}}
<pre>{{$var.Comment}}</pre>
<h4>{{$file.Path}}:{{$var.LineNo}} - {{$var.Source}}</h4>
{{end}}

{{end}}

</body>
`))

// Package site provides the docc.io website.
package site // import "docc.io/source/site"

import (
	"fmt"
	"html"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var Router = httprouter.New()

func init() {
	Router.GET("/github.com/:account/:name", gitHubGET)

	Router.GET("/source", sourceGET)
	Router.GET("/source/*path", sourceGET)
	Router.GET("/robots.txt", robotsTxtGET)

	Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notFound(w, `Nothing to see here; move along.`)
	})
}

func notFound(w http.ResponseWriter, msg string) {
	h := w.Header()
	h.Set("Content-Type", "text/html;charset=utf-8")
	h.Set("Content-Language", "en")
	w.WriteHeader(http.StatusNotFound)

	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html;charset=utf-8">
<title>docc.io - not found</title>
</head>
<body>
<h1>HTTP 404 &mdash; Not Found</h1>
<h4>%s</h4>
</body>
`, html.EscapeString(msg))
}

// sourceGet handles the Go canonical import path redirection.
func sourceGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h := w.Header()
	h.Set("Content-Type", "text/html;charset=utf-8")
	h.Set("Content-Language", "en")

	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html;charset=utf-8">
<meta name="go-import" content="docc.io/source git https://github.com/pascaldekloe/docc">
<meta name="go-source" content="docc.io/source https://github.com/pascaldekloe/docc/ https://github.com/pascaldekloe/docc/tree/master{/dir} https://github.com/pascaldekloe/docc/blob/master{/dir}/{file}#L{line}">
<title>docc.io - Go cananonical import path</title>
</head>
<body>
Nothing to see here; <a href="https://godoc.org/docc.io/source">move along</a>.
</body>
`))
}

func robotsTxtGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h := w.Header()
	h.Set("Content-Type", "text/plain;charset=utf-8")
	h.Set("Content-Language", "en")

	w.Write([]byte(`User-agent: *
Disallow: /source
`))
}

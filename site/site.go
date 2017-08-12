// Package site provides the docc.io website.
package site // import "docc.io/source/site"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var Router = httprouter.New()

func init() {
	Router.GET("/source", sourceGET)
	Router.GET("/source/*path", sourceGET)
	Router.GET("/robots.txt", robotsTxtGET)
}

func sourceGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h := w.Header()
	h.Set("Content-Type", "text/html;charset=UTF-8")
	h.Set("Content-Language", "en")

	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html;charset=UTF-8"/>
<meta name="go-import" content="docc.io/source git https://github.com/pascaldekloe/docc">
<meta name="go-source" content="docc.io/source https://github.com/pascaldekloe/docc/ https://github.com/pascaldekloe/docc/tree/master{/dir} https://github.com/pascaldekloe/docc/blob/master{/dir}/{file}#L{line}">
<title>docc.io - Go cananonical import path</title>
</head>
<body>
Nothing to see here; <a href="https://godoc.org/docc.io/source">move along</a>.
</body>
</html>
`))
}

func robotsTxtGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h := w.Header()
	h.Set("Content-Type", "text/plain;charset=UTF-8")
	h.Set("Content-Language", "en")

	w.Write([]byte(`User-agent: *
Disallow: /source
`))
}

// Package web provides the handlers for web display.
// Two types of Handlers :
//   - HTMX handlers : run computations, actual functions
//   - Web handlers : only display HTML, alias to variables.
package web

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:generate templ generate
//go:generate npm run build:css

//go:embed assets/js assets/css/output.css
var Assets embed.FS

func RegisterStaticRoutes(r chi.Router) {
	fileServer := http.FileServer(http.FS(Assets))

	r.Handle("/assets/*", fileServer)
}

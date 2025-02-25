package netutils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ServeRobotsTxt(r chi.Router, txt string) {
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(txt))
	})
}

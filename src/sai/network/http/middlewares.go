package http

import (
	"fmt"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/fatih/color"
	"net/http"
	"time"
)

func NoteTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		d := color.New(color.FgCyan, color.Bold)
		d.Println(fmt.Sprintf("%s\t%s\t%s\t%s\t%s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI, r.URL.Path[1:], time.Since(start)))
	})
}

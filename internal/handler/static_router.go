package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterWebRouter(server *rest.Server) {

	prefix := ""
	if e := os.Getenv("UI_PREFIX"); e != "" {
		if !strings.HasPrefix(e, "/") {
			e = "/" + e
		}
		prefix = e
	}

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: distHandler(),
			},
			{
				Method:  http.MethodGet,
				Path:    "/favicon.ico",
				Handler: distHandler(),
			},
			{
				Method:  http.MethodGet,
				Path:    fmt.Sprintf("%s/:1", prefix),
				Handler: stripPrefixHandler(prefix, "web/out"),
			},
			{
				Method:  http.MethodGet,
				Path:    fmt.Sprintf("%s/:1/:2", prefix),
				Handler: stripPrefixHandler(prefix, "web/out"),
			},
			{
				Method:  http.MethodGet,
				Path:    fmt.Sprintf("%s/:1/:2/:3/:4", prefix),
				Handler: stripPrefixHandler(prefix, "web/out"),
			},
		},
	)
}

func distHandler() func(http.ResponseWriter, *http.Request) {
	fs := http.FileServer(http.Dir("web/out"))
	return func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}
}

func stripPrefixHandler(pattern, fileDir string) http.HandlerFunc {
	if pattern == "" {
		pattern = "/"
	}
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(pattern, http.FileServer(http.Dir(fileDir)))
		handler.ServeHTTP(w, req)
	}
}

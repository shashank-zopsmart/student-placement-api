package middlewares

import "net/http"

func Middleware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !(req.Header.Get("x-api-key") == "a601e44e306e430f8dde987f65844f05" ||
			req.Header.Get("x-api-key") == "84dcb7c09b4a4af8a67f4577ffe9b255") {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
			return
		}

		if req.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Header Content-Type incorrect"))
			return
		}
		originalHandler.ServeHTTP(w, req)
	})
}

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// -------------------------------------------------
func main() {

	params := CollectParams()
	httpCli := GetHttpCli()
	target, _ := url.Parse(params.k8sEndpoint)

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// configure and enable oidc - if set
		if params.oidcEndpoint != "" {
			verifier := CreateOidcVerifier()
			authHeader := r.Header.Get("Authorization")
			rawToken := strings.TrimPrefix(authHeader, "Bearer ")
			if rawToken == "" {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			_, err := verifier.Verify(context.Background(), rawToken)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// forward request
		req, _ := http.NewRequest(r.Method, target.String()+r.URL.Path, r.Body)
		req.Header = r.Header.Clone()
		req.Header.Set("Authorization", "Bearer "+params.k8sToken)
		resp, err := httpCli.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// runner
	log.Println("Proxy runs on :8080")
	http.ListenAndServe(":8080", proxy)
}

package main

import (
	"context"
	"fmt"
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

		// allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Preflight-Request direkt beantworten
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// issue token for token auth
		if r.URL.Path == "/token" {
			baUsername, baPassword, ok := r.BasicAuth()
			fmt.Println(baUsername, baPassword)
			if baUsername != params.tokenAuthUsername || baPassword != params.tokenAuthPassword || !ok {
				fmt.Println("invalid credentials")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			jwt := IssueBearerToken(baUsername)
			w.Write([]byte(jwt))
			return
		}

		// procedure for jwt and oidc
		if !params.disableAuth {
			var authOk bool = false

			authHeader := r.Header.Get("Authorization")
			rawToken := strings.TrimPrefix(authHeader, "Bearer ")
			if rawToken == "" && !params.disableAuth {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			if params.tokenAuthUsername != "" && params.tokenAuthPassword != "" {
				authOk = CheckJwtToken(rawToken)
			}

			if !authOk && params.oidcEndpoint != "" {
				verifier := CreateOidcVerifier()
				_, err := verifier.Verify(context.Background(), rawToken)
				if err == nil {
					authOk = true
				}
			}

			if !authOk {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// forward request
		tgtUrl := target.String() + r.URL.Path
		log.Printf("%s -> %s", r.Method, tgtUrl)

		req, _ := http.NewRequest(r.Method, tgtUrl, r.Body)
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
	paramsInfo := params
	paramsInfo.k8sToken = ""
	log.Printf("Configured Params: %+v", paramsInfo)

	http.ListenAndServe(":8080", proxy)
}

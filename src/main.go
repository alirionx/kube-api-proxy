package main

import (
	"context"
	"encoding/base64"
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
		// configure and enable basic auth - if set
		if params.basicAuthUsername != "" && params.basicAuthPassword != "" {
			authHeader := r.Header.Get("Authorization")
			encoded := strings.TrimPrefix(authHeader, "Basic ")
			decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				fmt.Println("failed to decode:", err)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			decoded := string(decodedBytes)
			parts := strings.Split(decoded, ":")
			if len(parts) != 2 {
				fmt.Println("wrong basic auth format:", err)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if parts[0] != params.basicAuthUsername && parts[1] != params.basicAuthUsername {
				fmt.Println("invalid credentials:")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// configure and enable oidc - if set
		if params.oidcEndpoint != "" && params.basicAuthUsername == "" {
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

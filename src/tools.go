package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
)

// -------------------------------------------------
type Params struct {
	k8sToken              string
	k8sEndpoint           string
	basicAuthUsername     string
	basicAuthPassword     string
	oidcEndpoint          string
	oidcClientID          string
	containerK8sTokenPath string
}

// -------------------------------------------------
func CollectParams() Params {
	params := Params{
		containerK8sTokenPath: "/var/run/secrets/kubernetes.io/serviceaccount/token",
	}

	// gether relevant envs----------------
	k8s_ep := os.Getenv("K8S_ENDPOINT")
	k8s_sh := os.Getenv("KUBERNETES_SERVICE_HOST")
	k8s_sp := os.Getenv("KUBERNETES_SERVICE_PORT")
	k8s_tk := os.Getenv("K8S_TOKEN")

	params.basicAuthUsername = os.Getenv("BASIC_AUTH_USERNAME")
	params.basicAuthPassword = os.Getenv("BASIC_AUTH_PASSWORD")

	params.oidcEndpoint = os.Getenv("OIDC_ENDPOINT")
	params.oidcClientID = os.Getenv("OIDC_CLIENT_ID")

	// determinate K8S API Endpoint--------
	if k8s_ep != "" {
		params.k8sEndpoint = k8s_ep
	} else if k8s_sh != "" && k8s_sp != "" {
		params.k8sEndpoint = fmt.Sprintf("https://%s:%s", k8s_sh, k8s_sp)
	} else {
		log.Fatal("unable to determinate k8s endpoint")
	}

	// determinate K8S API TOKEN-----------
	if k8s_tk != "" {
		params.k8sToken = k8s_tk
	} else if _, err := os.Stat(params.containerK8sTokenPath); err == nil {
		data, err := os.ReadFile(params.containerK8sTokenPath)
		if err != nil {
			log.Fatal("unable to read service account token", err)
		}
		params.k8sToken = string(data)
	} else {
		log.Fatal("unable to determinate k8s token")
	}

	// OIDC CONFIG Check-------------------
	if params.oidcEndpoint != "" {
		_, err := url.Parse(params.oidcEndpoint)
		if err != nil {
			log.Fatal("invalid oidc endpoint", err)
		}
	}

	return params
}

// ------------
func GetHttpCli() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return *client
}

// ------------
func CreateOidcVerifier() oidc.IDTokenVerifier {
	params := CollectParams()
	provider, err := oidc.NewProvider(context.Background(), params.oidcEndpoint)
	if err != nil {
		log.Fatal("oidc provider error:", err)
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID:          params.oidcClientID,
		SkipClientIDCheck: true,
	})
	return *verifier
}

### Envs and things ;)
```
export K8S_ENDPOINT="https://192.168.10.10:6443"
export K8S_TOKEN="eyJ..."

export BASIC_AUTH_USERNAME="admin"
export BASIC_AUTH_PASSWORD="admin"
export OIDC_ENDPOINT=""
export OIDC_CLIENT_ID=""

export OIDC_ENDPOINT="https://keycloak.app-scape.de/realms/master"
export OIDC_CLIENT_ID="appscape"
export BASIC_AUTH_USERNAME=""
export BASIC_AUTH_PASSWORD=""


docker build -t ghcr.io/alirionx/kube-api-proxy:latest --no-cache .
docker run \
  -d \
  -rm \
  -p 8080:8080 \
  -e K8S_ENDPOINT="https://192.168.10.10:6443" \
  -e K8S_TOKEN="ey..." \
  -e OIDC_ENDPOINT="https://keycloak.app-scape.de/realms/master" \
  -e OIDC_CLIENT_ID="appscape" \
  ghcr.io/alirionx/kube-api-proxy:latest

```
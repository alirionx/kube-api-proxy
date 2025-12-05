## A simple microservice for k8s api proxying with optional OIDC or basic auth
- written in golang
- uses the default net library 

<br>

__Notes:__  
This microservice enables separate exposing of the Kubernetes API of a cluster.  
Depending on the permissions of the respective Service Accounts or Kubernetes users (via Kubeconfig), access rights can be managed.  
It is possible to secure the exposing through OIDC or Basic Authentication.  
SSL termination can be handled, for example, via an Ingress Controller when the microservice is running inside a Kubernetes cluster.


<br>


__Config Envs__  
- K8S_ENDPOINT 
- K8S_TOKEN
- BASIC_AUTH_USERNAME
- BASIC_AUTH_PASSWORD
- OIDC_ENDPOINT
- OIDC_CLIENT_ID

<br>

__Docker Example Usage:__  
```
docker build -t kube-api-proxy:latest --no-cache .
docker run \
  -d \
  -rm \
  -p 8080:8080 \
  -e K8S_ENDPOINT="https://192.168.10.10:6443" \
  -e K8S_TOKEN="ey..." \
  -e OIDC_ENDPOINT="https://keycloak.example.com/realms/apps" \
  -e OIDC_CLIENT_ID="apps" \
  kube-api-proxy:latest
```

<br>

__K8s Example Usage:__  
```
kubectl apply -f k8s-manifest 
```
Do not forget to set your individual parameter in the manifest like e.g:
- Cluster Role Rules
- Envs for K8S Endpoint, OIDC Backend or BASIC Auth


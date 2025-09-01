# quote-of-the-day
Project with kubernetes

docker build -t quote-api .
docker run -p 8080:8080 --rm quote-api

docker tag quote-api your-username/quote-api:1.0

kubectl apply -f quote-configmap.yaml

kubectl get configmap quote-app-config -o yaml

quote-app-deployment-64488f44cd-cmdkt no

# Record HPA results
1. kubectl get hpa -n quote-app-ns -w
2. kubectl get pods -n quote-app-ns -w
3. kubectl port-forward svc/quote-app-service 8080:80 -n quote-app-ns
4. Write-Host "Generating load... Press Ctrl+C to stop."; while ($true) { curl.exe -s http://localhost:8080/quote | Out-Null }

## Results that should be observed
In the HPA watcher, you'll see the TARGETS percentage climb above 50%.
Once the threshold is breached, the REPLICAS count will jump from 2 to 3.
In the Pod watcher, you'll see a new pod appear with the status ContainerCreating, which will soon change to Running.
If you keep the load running, the HPA will continue adding pods until it reaches your maxReplicas limit of 5.

# Ingress Controller
An Ingress controller is a specialized reverse proxy that runs in your cluster. It acts like a smart traffic manager or a receptionist for your applications.
It works with two key components:
    - The Ingress Controller: The actual software (in our case, NGINX) that runs in a pod, listening for external traffic.
    - The Ingress Resource: A Kubernetes YAML file you create that defines the routing rules (e.g., "send traffic for api.example.com/quote to the quote-app-service").
    - The controller watches for Ingress resources and automatically configures itself to follow the rules you define.
### Install
```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/cloud/deploy.yaml
```

### Verify
```
kubectl get pods -n ingress-nginx -w
```
### Details
spec.rules: This section defines the core logic. The rule we've created says:
host: "quotes.local": When traffic arrives with the hostname quotes.local...
path: /: ...for any URL path (/ with a Prefix match type)...
backend.service: ...route the traffic to the quote-app-service on its port 80.
The overall traffic flow will be: External Request -> Ingress Controller (NGINX) -> Ingress Rule -> quote-app-service -> One of the Pods.

## Clean install
```
helm uninstall quote-release --namespace quote-app-ns
kubectl create namespace quote-app-ns
helm install quote-release ./quote-chart --namespace quote-app-ns
```

# Using Argo Rollouts with many microservices

Example with Argo Rollouts with backend and frontend

![Dashboard](src/loan/static/diagram.svg)

## Run only the backend manually

Install [GoLang](https://go.dev/) locally

```
cd src/interest
APP_VERSION=1.2  go run .
```
You can now acces the backend at `http://localhost:8080`

## Run only the frontend manually

Install [GoLang](https://go.dev/) locally

```
cd src/loan
APP_VERSION=1.4 PORT=9000 BACKEND_HOST=localhost go run .
```

You can now access the backend at `http://localhost:9000`

## Run both of them at the same time

Install [Docker compose](https://docs.docker.com/compose/) (no need for local GoLang installation)

```
cd src
docker compose up
```

And now you can use the same URLs as above to access the services.

## Run on Kubernetes as deployments

```
cd manifests/plain
kubectl create ns plain
kubectl apply -f . -n plain
kubectl port-forward svc/my-plain-backend-service 8000:8080 -n plain
kubectl port-forward svc/my-plain-frontend-service 9000:8080 -n plain
```

You can now access the backend at `http://localhost:8000` and the backend at `http://localhost:9000`
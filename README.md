# multi-service-argo-rollouts-example
Example with Argo Rollouts with many services

## Run only the backend manually

Install [GoLang](https://go.dev/) locally

```
cd interest
APP_VERSION=1.2  go run .
```
You can now acces the backend at `http://localhost:8080`

## Run only the frontend manually

Install [GoLang](https://go.dev/) locally

```
cd loan
APP_VERSION=1.4 PORT=9000 BACKEND_HOST=localhost go run .
```

You can now acces the backend at `http://localhost:9000`

## Run both of them at the same time

Install [Docker compose](https://docs.docker.com/compose/) (no need for local GoLang installation)

```
docker compose up
```

And now you can use the same URLs as above to access the services.

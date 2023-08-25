# multi-service-argo-rollouts-example
Example with Argo Rollouts with many services

## Run the backend

```
cd interest
APP_VERSION=1.2  go run .
```
You can now acces the backend at `http://localhost:8080`

## Run the frontend

```
cd loan
APP_VERSION=1.4 PORT=9000 BACKEND_HOST=localhost go run .
```

You can now acces the backend at `http://localhost:9000`

## Run both of them

```
docker compose up
```

And now you the same URLs as above to access the services

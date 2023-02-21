#  Url info checker

```
Url info checker
```

## Development
```
Run go run cmd/server/main.go -env=env in cmd/server

You need to run local postgresql to be able to use APIs
Look db/README.md for migrations and local postgresql db
```

## Testing

```
Run  go test ./... in root folder
```

## Swagger

```
Open http://127.0.0.1:9234/static/swagger-ui/ for swagger UI
```

## Using Docker to simplify development (optional)

```
Run docker-compose run api migrate -path /api/migrations/ -database "postgres://user:password@postgres:5432/dbname?sslmode=disable" up
 in root folder
Run docker-compose up 
in root folder
```

## Metrics 

```
Open http://localhost:9090/ for Prometheus
Open http://localhost:16686/search for Jaeger UI
```

## Vault 

```
You can use vault as a secure secrets management but you should add _SECURE  to environment variables like 
 DATABASE_USERNAME_SECURE="/database:username"
 DATABASE_PASSWORD_SECURE="/database:password"
 You can access vault interface at localhost:8300 
 Method Token 
 Token = myroot
```




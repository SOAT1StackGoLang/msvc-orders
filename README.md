# msvc-orders

Microservice Orders

## Owasp Scanning Report

[Owasp Report](./owasp/README.md)

## Swagger Documentation

[Swagger Documentation](./docs/readme.md)

## Generate Swagger

```bash
go get -u github.com/swaggo/swag/cmd/swag
go install github.com/swaggo/swag/cmd/swag
cd internal/transport/routes
swag fmt
swag init -g helpers.go -o ../../../docs/
```

- if you not find the swag binary, check if your go PATH is on your path or use ~/go/bin/swag

More info on the formats at [Swag](https://github.com/swaggo/swag?tab=readme-ov-file#api-operation) offical

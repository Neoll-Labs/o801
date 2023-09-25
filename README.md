

Open API documentation
[OpenAPI documentation swagger](https://editor-next.swagger.io/?url=https://raw.githubusercontent.com/nelsonstr/o801/main/api/openapi.yaml)

## interfaces definition

* Handlers interface definition for creating and retrieving [Handlers](./internal/interfaces/handlers.go)
* Service interface definition for creating and retrieving [Service](./internal/interfaces/service.go)
* Cache interface definition for creating and retrieving [Cache](./internal/interfaces/cache.go)
* Repository interface definition for creating and retrieving [Repository](./internal/interfaces/repository.go)

## Models

There are two models for User:
* A [store model](./internal/models/user.go)
* A [view model](./internal/models/user.go)


[deepsource code](https://app.deepsource.com/gh/nelsonstr/o801)

* Docker compose Services
  * backend
  * DB PostgreSQL


### command line

#### docker build container


```shell
  docker build . --target=prod --tag=o801:latest
```

#### Starting docker compose

* Run the following command to start docker compose

```shell
docker-compose up
```

* validate openapi.yaml
```shell
docker run --rm -v C:/data/o801:/local openapitools/openapi-generator-cli:v7.0.0 validate -i /local/api/openapi.yaml
```

* create user
```shell
curl 127.0.0.1:8080/api/v1/users/ -X POST -H 'accept: application/json' -H 'Content-Type: application/json' -d'{"name": "nelson"}'
```

* get user
```shell
curl 127.0.0.1:8080/api/v1/users/1
```

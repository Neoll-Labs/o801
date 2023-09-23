

* Docker compose Services
  * backend
  * DB PostgreSQL

* docker build

```shell
  docker build . --target=prod --tag=o801:latest
```

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
curl 127.0.0.1:8080/api/v1/users/ -X POST -H 'accept: application/json' -H 'Content-Type: application/json' -d'{"Name": "nelson"}'
```

* get user
```shell
curl 127.0.0.1:8080/api/v1/users/1
```

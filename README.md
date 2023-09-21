
* docker build

```shell
  docker build . --target=prod --tag=o801:latest
```

* validate openapi.yaml
```shell
docker run --rm -v C:/data/o801:/local openapitools/openapi-generator-cli:v7.0.0 validate -i /local/api/openapi.yaml
```


```shell

curl 127.0.0.1:8080/api/v1/auth/login
curl 127.0.0.1:8080/api/v1/auth/login -X POST
curl 127.0.0.1:8080/api/v1/auth/logout

curl 127.0.0.1:8080/api/v1/users/ -X POST

curl http://127.0.0.1:8080/metrics/

```

```shell

 curl 127.0.0.1:8080/api/v1/users/ -X POST\
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "Name": "nelson"
}'
```
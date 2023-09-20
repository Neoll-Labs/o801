
* docker build

```shell
  docker build . --target=prod --tag=o801:latest
```

* validate openapi.yaml
```shell
docker run --rm -v C:/data/o801:/local openapitools/openapi-generator-cli:v7.0.0 validate -i /local/api/openapi.yaml
```

* rebuild server
```shell

docker run --rm -v C:/data/o801:/local openapitools/openapi-generator-cli:v7.0.0 generate -i /local/api/openapi.yaml -g go-server -o /local/
```


* swagger ui
```html
http://127.0.0.1:8080/swagger/index.html
```

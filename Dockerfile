# STAGE 1: building the executable
FROM golang:1.21-alpine AS build

WORKDIR /src

COPY ./go.mod ./go.sum ./

# download dependencies
RUN go mod download

COPY ./ ./

# Run tests
RUN CGO_ENABLED=0 go test -timeout 30s -v ./... && \
    CGO_ENABLED=0 go build -o /app ./cmd

# STAGE 2: build the container to run distroyless non root
FROM gcr.io/distroless/static@sha256:92d40eea0b5307a94f2ebee3e94095e704015fb41e35fc1fcbd1d151cc282222 AS prod


# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app /app

EXPOSE 8080
# run binary; use vector form
ENTRYPOINT ["/app"]

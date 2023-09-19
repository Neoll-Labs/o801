ARG GO_VERSION=1.21

# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache git

WORKDIR /src

COPY ./go.mod ./go.sum ./

# download dependencies
RUN go mod download


COPY ./ ./

# Run tests
# RUN CGO_ENABLED=0 go test -timeout 30s -v github.com/nelsonstr/o801

# Build the executable
RUN CGO_ENABLED=0 go build -o /app ./cmd


# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS prod
USER nonroot:nonroot

# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app /app

# run binary; use vector form
ENTRYPOINT ["/app"]
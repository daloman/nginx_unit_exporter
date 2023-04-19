FROM golang:1.19.5-alpine3.17 AS build_deps

RUN apk add --no-cache git

WORKDIR /workspace

COPY go.mod .

COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o nginx_unit_exporter -ldflags '-w -extldflags "-static"' .

FROM alpine:3.17

RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/nginx_unit_exporter /usr/local/bin/nginx_unit_exporter

ENTRYPOINT ["nginx_unit_exporter"]

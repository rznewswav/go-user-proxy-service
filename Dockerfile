FROM golang:1.19.4-alpine3.17 AS build
RUN apk update && apk add gcc libc-dev build-base
WORKDIR /usr/local/go
# copied from https://www.docker.com/blog/containerize-your-go-developer-environment-part-2/
ENV CGO_ENABLED=1
COPY go.* .
RUN go mod download
COPY . .
RUN go build -tags musl -o build.exe cmd/main/main.go

FROM alpine:3.17 AS production
WORKDIR /usr/local/go
COPY --from=build /usr/local/go/build.exe /usr/local/go/start ./
COPY --from=build /usr/local/go/resources  ./resources
EXPOSE 3000

CMD ["/usr/local/go/build.exe"]
# uncomment once bugsnag has been integrated
# CMD ["/bin/sh", "start"]
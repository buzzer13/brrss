FROM golang:1.21-alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o brrss cmd/brrss-server/main.go

FROM alpine:3.19

EXPOSE 8080

RUN apk --no-cache add ca-certificates
COPY --from=build /app/brrss /usr/bin/brrss
USER nobody

CMD ["/usr/bin/brrss"]

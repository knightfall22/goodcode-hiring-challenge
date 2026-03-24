FROM golang:1.24-alpine AS builder
RUN mkdir -p /go/src/api
WORKDIR /go/src/api
COPY . .

RUN GIT_TERMINAL_PROMPT=1 \
    GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go build -v --installsuffix cgo --ldflags="-s" -o api ./cmd/server/main.go
FROM alpine:3.13

# convert build-arg to env variables
RUN apk add --no-cache tzdata
ENV TZ=Europe/Bucharest
RUN mkdir -p /svc/
COPY --from=builder /go/src/api/api /svc/

EXPOSE 8484

WORKDIR /svc/

CMD ["./api"]

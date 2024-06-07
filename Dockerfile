#Build stage
FROM golang:1.22.2-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
#RUN apk --no-cache add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz    

#Run Stage
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/main .
#COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY flex-1656360315201-firebase-adminsdk-cgon8-b5cb17e022.json .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
FROM golang:1.21

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build cmd/scheduler.go

EXPOSE 8080

CMD ["/app/scheduler"]
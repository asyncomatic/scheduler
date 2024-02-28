FROM golang:1.21 as builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download && \
    go mod verify

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /dist/scheduler cmd/scheduler.go

FROM golang:1.21
LABEL org.opencontainers.image.source=https://github.com/asyncomatic/scheduler

COPY --from=builder /dist/scheduler .

USER 65534
EXPOSE 8080

ENTRYPOINT ["./scheduler"]
FROM golang:1.24.5-alpine3.22 AS build-stage

WORKDIR /marketplace

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY scripts/ ./scripts

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest

WORKDIR /marketplace

COPY --from=build-stage /marketplace/main .

EXPOSE 8080

CMD ["./main"]
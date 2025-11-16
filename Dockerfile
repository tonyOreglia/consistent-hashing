FROM golang:1.25.4

EXPOSE 8080

WORKDIR /usr/src/consistent_hash

COPY go.mod go.sum /usr/src/consistent_hash/
RUN go mod download

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

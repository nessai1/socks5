FROM golang:1.21 as builder

WORKDIR /socks5

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .

FROM scratch

COPY --from=builder /socks5/service .

CMD ["./service"]
FROM golang:1.17-alpine

WORKDIR /app
COPY go.mod go.sum logs.log ./

RUN go mod download

COPY . .

RUN go build -o /main ./cmd/compressor

EXPOSE 9001

CMD ["/main"]
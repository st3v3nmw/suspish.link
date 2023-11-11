From golang:1.21.3-alpine

ENV GIN_MODE=release

WORKDIR /go/src/app

COPY . .

RUN go build -o suspish ./...

CMD ["./suspish", "--verbose"]

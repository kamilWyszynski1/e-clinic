FROM golang:1.14

RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN ls

RUN go build ./cmd/main.go

CMD ["./main"]
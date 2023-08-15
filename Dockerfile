FROM golang:1.21

RUN mkdir -p /go/src/test-go-gin-gonic
WORKDIR /go/src/test-go-gin-gonic
COPY . .

RUN go mod download

RUN go build -o /main main.go

EXPOSE 8080

ENTRYPOINT ["/main"]
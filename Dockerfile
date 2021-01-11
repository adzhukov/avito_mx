FROM golang:1.15 AS builder

WORKDIR /src

ADD go.mod go.sum /src/
RUN go mod download

ADD . .
RUN go build -o /go/bin/avito_mx

ENTRYPOINT ["avito_mx"]

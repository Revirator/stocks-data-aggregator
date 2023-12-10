FROM golang:1.21.4-alpine3.18

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o /cfd

EXPOSE 8080

ENTRYPOINT ["/cfd"]
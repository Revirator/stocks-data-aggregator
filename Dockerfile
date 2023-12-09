FROM golang:1.21.4-alpine3.18

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY .env ./
COPY *.go ./

RUN go build -o /cfd

ADD templates ./templates
ADD static ./static

EXPOSE 8080

ENTRYPOINT ["/cfd"]
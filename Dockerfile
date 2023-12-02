FROM golang:1.21.4-alpine3.18

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY .env ./
COPY *.go ./

ADD templates /app/templates

RUN go build -o /stocks-data-aggregator

EXPOSE 8080

ENTRYPOINT ["/stocks-data-aggregator"]
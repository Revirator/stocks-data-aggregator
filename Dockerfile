FROM golang:1.21.4-alpine3.18

WORKDIR /app

COPY go.* ./

RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@v0.2.543

COPY . .

RUN templ generate
RUN go build -o /cfd

EXPOSE 8080

ENTRYPOINT ["/cfd"]
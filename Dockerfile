FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go build -o /main main.go

EXPOSE 8080

ENTRYPOINT [ "/main" ]
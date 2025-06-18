FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go build -o /main main.go

EXPOSE 3000

ENTRYPOINT [ "/main" ]
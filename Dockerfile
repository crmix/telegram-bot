FROM golang:1.23-alpine

WORKDIR /app

COPY . /app/

RUN go build -o main main.go
RUN chmod +x ./main

CMD ["./main"]
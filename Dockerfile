FROM golang:1.19

WORKDIR /app

COPY . .

RUN go build -o /go-docker-demo main/main.go

EXPOSE 8080

CMD [ "/go-docker-demo" ]
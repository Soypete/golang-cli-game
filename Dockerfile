FROM docker.io/golang:1.20-alpine

RUN apk update
RUN apk add git

WORKDIR /app
COPY go.* ./
RUN go mod download

EXPOSE 3000

COPY . ./
RUN go build -v game.go

CMD ["/app/game"]

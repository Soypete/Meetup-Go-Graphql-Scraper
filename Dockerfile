FROM golang:alpine

RUN apk update
RUN apk add git gcc g++

WORKDIR /app
ENV CGO_ENABLED=1
COPY . ./
RUN go mod download
RUN go mod tidy
RUN go build -v -o main.go

CMD ["/app/main"]

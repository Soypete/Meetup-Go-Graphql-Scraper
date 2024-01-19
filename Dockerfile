FROM golang:alpine

RUN apk update
RUN apk add git gcc g++

WORKDIR /app
ENV CGO_ENABLED=1
COPY . ./
RUN ls -la
RUN go mod download
RUN go mod tidy
RUN go build -v -tags=LD_LIBRARY_PATH=/duckdb -o main

CMD ["/app/main"]

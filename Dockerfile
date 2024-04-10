FROM golang:1.16.3-alpine3.13 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/sublinkX
EXPOSE 8000
CMD ["/app/sublinkX"]
FROM golang:1.21-alpine

WORKDIR /app

COPY . .
WORKDIR collector
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /collector

CMD [ "/collector" ]